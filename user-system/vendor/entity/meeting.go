package entity

import (
	"convention/codec"
	"convention/homererror"
	"io"
	"time"
	log "util/logger"
)

const TimeLayout = time.RFC3339

// Identifier
type MeetingTitle string

func (t MeetingTitle) Empty() bool {
	return t == ""
}
func (t MeetingTitle) Valid() bool {
	return !t.Empty() // NOTE: may not only !empty
}
func (t MeetingTitle) String() string {
	return string(t)
}

type MeetingInfo struct {
	Title         MeetingTitle
	Sponsor       *User
	Participators UserList

	StartTime time.Time
	EndTime   time.Time
}
type MeetingInfoSerializable struct {
	Title         MeetingTitle
	Sponsor       Username
	Participators []Username
	StartTime     string
	EndTime       string
}

type MeetingInfoForDatabase struct {
	Title         MeetingTitle `gorm:"primary_key"`
	Sponsor       Username
	Participators []UserInfoSerializable `gorm:"many2many:participations;"`
	// Participators []UserInfoSerializable `gorm:"many2many:meetingInfoForDataBase_userInfoSerializables;"`

	// TODO: with `gorm`, seems able to use `time.Time` ...
	StartTime time.Time
	EndTime   time.Time
}

type Meeting struct {
	MeetingInfo
}

func NewMeeting(info MeetingInfo) *Meeting {
	m := new(Meeting)
	m.MeetingInfo = info
	return m
}

func LoadMeeting(decoder codec.Decoder, m *Meeting) {
	mInfoSerial := new(MeetingInfoSerializable)

	err := decoder.Decode(mInfoSerial)
	if err != nil {
		log.Fatal(err)
	}
	m.MeetingInfo = *(mInfoSerial.Deserialize())
}
func LoadedMeeting(decoder codec.Decoder) *Meeting {
	m := new(Meeting)
	LoadMeeting(decoder, m)
	return m
}

func (m *Meeting) Save(encoder codec.Encoder) error {
	return encoder.Encode(*m.MeetingInfo.Serialize())
}

func (info *MeetingInfo) Serialize() *MeetingInfoSerializable {
	if info == nil {
		log.Fatal("nil *MeetingInfo calls Seriablize ...\n")
	}
	mInfoSerial := new(MeetingInfoSerializable)

	mInfoSerial.Title = info.Title

	if sponsor := info.Sponsor; sponsor != nil {
		mInfoSerial.Sponsor = sponsor.ID
	}

	mInfoSerial.Participators = info.Participators.Identifiers()

	mInfoSerial.StartTime = info.StartTime.Format(TimeLayout)
	mInfoSerial.EndTime = info.EndTime.Format(TimeLayout)

	return mInfoSerial
}
func (infoSerial *MeetingInfoSerializable) Deserialize() *MeetingInfo {
	if infoSerial == nil {
		log.Fatal("nil *MeetingInfoSerializable calls Deseriablize ...\n")
	}
	info := new(MeetingInfo)

	info.Title = infoSerial.Title

	// CHECK: Need ensure Sponsor not nil ?
	info.Sponsor = infoSerial.Sponsor.RefInAllUsers()

	info.Participators.InitFrom(infoSerial.Participators)

	// NOTE: better code ?
	var err1, err2 error
	info.StartTime, err1 = time.Parse(TimeLayout, infoSerial.StartTime)
	info.EndTime, err2 = time.Parse(TimeLayout, infoSerial.EndTime)
	if err1 != nil || err2 != nil {
		log.Fatalf("time.Parse fail when parsing %v / %v\n", infoSerial.StartTime, infoSerial.EndTime)
	}

	return info
}

// ................................................................

type MeetingList struct {
	Meetings map[MeetingTitle](*Meeting)
}

type MeetingListRaw = []*Meeting

type MeetingInfoSerializableList []MeetingInfoSerializable

type MeetingInfoListPrintable = MeetingInfoSerializableList

func (ml *MeetingList) Serialize() MeetingInfoSerializableList {
	ret := make(MeetingInfoSerializableList, 0, ml.Size())

	ml.ForEach(func(m *Meeting) error {
		if m == nil {
			log.Warning("A nil Meeting is to be used. Just SKIP OVER it.")
			return nil
		}
		ret = append(ret, *(m.MeetingInfo.Serialize()))
		return nil
	})

	return ret
}

func (mlSerial MeetingInfoSerializableList) Size() int {
	return len(mlSerial)
}

func (mlSerial MeetingInfoSerializableList) Deserialize() *MeetingList {
	ret := NewMeetingList()

	for _, mInfoSerial := range mlSerial {
		if mInfoSerial.Title.Empty() {
			log.Warning("A No-Title MeetingInfo is to be used. Just SKIP OVER it.")
			continue
		}

		m := NewMeeting(*(mInfoSerial.Deserialize()))
		if err := ret.Add(m); err != nil {
			log.Error(err)
		}
	}
	return ret
}

func NewMeetingList() *MeetingList {
	ml := new(MeetingList)
	ml.Meetings = make(map[MeetingTitle](*Meeting))
	return ml
}

// CHECK: Need in-place load method ?

func LoadMeetingList(decoder codec.Decoder, ml *MeetingList) {
	// CHECK: Need clear ml ?

	mlSerial := new(MeetingInfoSerializableList)
	if err := decoder.Decode(mlSerial); err != nil {
		switch err {
		case io.EOF:
			// FIXME: not sure io.EOF would always indicate empty Decoder, however I don't think this check should be placed otherwhere
			break
		default:
			log.Fatal(err)
		}
	}
	for _, mInfoSerial := range *mlSerial {
		m := NewMeeting(*(mInfoSerial.Deserialize()))
		if err := ml.Add(m); err != nil {
			log.Error(err)
		}
	}
}
func (ml *MeetingList) LoadFrom(decoder codec.Decoder) {
	LoadMeetingList(decoder, ml)
}

func LoadedMeetingList(decoder codec.Decoder) *MeetingList {
	ml := NewMeetingList()
	LoadMeetingList(decoder, ml)
	return ml
}

func (ml *MeetingList) Identifiers() []MeetingTitle {
	ret := make([]MeetingTitle, 0, ml.Size())
	for _, mInfoSerial := range ml.Textualize() {
		ret = append(ret, mInfoSerial.Title)
	}
	return ret
}

// CHECK: better name or conduct ? `ul.PublicInfos` <---> `ml.Textualize`
func (ml *MeetingList) Textualize() MeetingInfoListPrintable {
	return ml.Serialize()
}

func (ml *MeetingList) Save(encoder codec.Encoder) error {
	sl := ml.Serialize()
	return encoder.Encode(sl)
}

func (ml *MeetingList) Size() int {
	return len(ml.Meetings)
}

func (ml *MeetingList) Ref(title MeetingTitle) *Meeting {
	return ml.Meetings[title] // NOTE: if directly return accessed result from a map like this, would not get the (automatical) `ok`
}
func (ml *MeetingList) Contains(title MeetingTitle) bool {
	m := ml.Ref(title)
	return m != nil
}

func (ml *MeetingList) Add(meeting *Meeting) error {
	if meeting == nil {
		return homererror.ErrNilMeeting
	}
	title := meeting.Title
	if ml.Contains(title) {
		return homererror.ErrExistedMeeting
	}
	ml.Meetings[title] = meeting
	return nil
}
func (ml *MeetingList) Remove(meeting *Meeting) error {
	if meeting == nil {
		return homererror.ErrNilMeeting
	}
	title := meeting.Title
	if ml.Contains(title) {
		delete(ml.Meetings, title) // NOTE: never error, according to 'go-maps-in-action'
		return nil
	}
	return homererror.ErrMeetingNotFound
}
func (ml *MeetingList) PickOut(title MeetingTitle) (*Meeting, error) {
	if title.Empty() {
		return nil, homererror.ErrEmptyMeetingTitle
	}
	m := ml.Ref(title)
	if m == nil {
		return nil, homererror.ErrMeetingNotFound
	}
	defer ml.Remove(m)
	return m, nil
}

func (ml *MeetingList) Slice() MeetingListRaw {
	meetings := make(MeetingListRaw, 0, ml.Size())
	for _, m := range ml.Meetings {
		meetings = append(meetings, m) // CHECK: maybe better to use index in golang ?
	}
	return meetings
}

// ForEach used to extension/concrete logic for whole MeetingList
func (ml *MeetingList) ForEach(fn func(*Meeting) error) error {
	for _, v := range ml.Meetings {
		if err := fn(v); err != nil {
			// CHECK: Or, lazy error ?
			return err
		}
	}
	return nil
}

// Filter used for all extension/concrete select for whole MeetingList
func (ml *MeetingList) Filter(pred func(Meeting) bool) *MeetingList {
	ret := NewMeetingList()
	for _, m := range ml.Meetings {
		if pred(*m) {
			ret.Add(m)
		}
	}
	return ret
}
