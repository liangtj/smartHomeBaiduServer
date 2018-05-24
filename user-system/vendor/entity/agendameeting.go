package entity

import (
	"convention/codec"
	"convention/homererror"
	log "util/logger"
)

var allMeetings *MeetingList

func (title MeetingTitle) RefInAllMeetings() *Meeting {
	return allMeetings.Ref(title)
}
func GetAllMeetings() *MeetingList {
	return allMeetings
}

// NOTE: TODEL: No need now
// var dissolvedMeetings = *(NewMeetingList())

// func (title MeetingTitle) RefInDissolvedMeetings() *Meeting {
// 	return dissolvedMeetings.Ref(title)
// }
// func GetDissolvedMeetings() *MeetingList {
// 	return &dissolvedMeetings
// }

// LoadAllMeeting concretely loads all Meetings
func LoadAllMeeting(decoder codec.Decoder) {
	allMeetings = LoadedMeetingList(decoder)
}

// SaveAllMeeting concretely saves all Meetings
func SaveAllMeeting(encoder codec.Encoder) error {
	meetings := GetAllMeetings()
	return meetings.Save(encoder)
}

//

// SponsoredBy checks if Meeting sponsored by User
func (m *Meeting) SponsoredBy(name Username) bool {
	if m == nil {
		return false
	}
	if m.Sponsor == nil {
		log.Warningln("m.SponsoredBy(name) where m.Sponsor == nil.\n")
		return false
	}
	return m.Sponsor.ID == name
}

// ContainsParticipator checks if Meeting's participators contains the User
func (m *Meeting) ContainsParticipator(name Username) bool {
	if m == nil {
		return false
	}
	return m.Participators.Contains(name)
}

// Dissolve deletes the Meeting (, not by a User)
func (m *Meeting) Dissolve() error {
	if m == nil {
		return homererror.ErrNilMeeting
	}

	if err := GetAllMeetings().Remove(m); err != nil {
		log.Print(err)
		return err
	}
	log.Printf("Meeting %v is dissolved.\n", m.Title)
	return nil
}

// Exclude removes User from Meeting's participators list
func (m *Meeting) Exclude(u *User) error {
	if m == nil {
		return homererror.ErrNilMeeting
	}
	if u == nil {
		return homererror.ErrNilUser
	}

	if err := m.Participators.Remove(u); err != nil {
		log.Print(err)
		return err
	}
	if m.Participators.Size() <= 0 {
		return m.Dissolve()
	}
	log.Printf("User %v is excluded from Meeting %v.\n", u.ID, m.Title)
	return nil
}

// Involve adds User to Meeting's participators list
func (m *Meeting) Involve(u *User) error {
	if m == nil {
		return homererror.ErrNilMeeting
	}
	if u == nil {
		return homererror.ErrNilUser
	}

	if err := m.Participators.Add(u); err != nil {
		log.Print(err)
		return err
	}
	log.Printf("User %v is involved in Meeting %v.\n", u.ID, m.Title)
	return nil
}
