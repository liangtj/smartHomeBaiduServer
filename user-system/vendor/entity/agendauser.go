package entity

import (
	"convention/codec"
	"convention/homererror"
	"time"
	log "util/logger"
)

// CHECK: Not sure where to place ...
var allUsersRegistered *UserList

// RefInAllUsers returns the ref of a Registered User depending on the Username
func (name Username) RefInAllUsers() *User {
	return allUsersRegistered.Ref(name)
}

// GetAllUsersRegistered returns the reference of the UserList of all Registered Users
func GetAllUsersRegistered() *UserList {
	return allUsersRegistered
}

// LoadUsersAllRegistered concretely loads all Registered Users
func LoadUsersAllRegistered(decoder codec.Decoder) {
	allUsersRegistered = LoadedUserList(decoder)
}

// SaveUsersAllRegistered concretely saves all Registered Users
func SaveUsersAllRegistered(encoder codec.Encoder) error {
	users := allUsersRegistered
	return users.Save(encoder)
}

// ........

func (u *User) Registered() bool {
	if u == nil {
		return false
	}
	return GetAllUsersRegistered().Contains(u.Name)
}

func (u *User) involvedMeetings() *MeetingList {
	return GetAllMeetings().Filter(func(m Meeting) bool {
		return m.SponsoredBy(u.Name) || m.ContainsParticipator(u.Name)
	})
}

func (u *User) FreeWhen(start, end time.Time) bool {
	if u == nil {
		return false
	}

	// NOTE: need improve:
	if err := u.involvedMeetings().ForEach(func(m *Meeting) error {
		s1, e1 := m.StartTime, m.EndTime
		s2, e2 := start, end
		if s1.Before(e2) && e1.After(s2) {
			return homererror.ErrConflictedTimeInterval
		}
		return nil
	}); err != nil {
		log.Error(err)
		return false
	}

	return true
}

// QueryAccount queries an account, where User as the actor
func (u *User) QueryAccount() error {
	return homererror.ErrNeedImplement
}

// QueryAccountAll queries all accounts, where User as the actor
func (u *User) QueryAccountAll() []UserInfoPublic {
	// NOTE: whatever, temporarily ignore the problem that the actor of query is Nil
	username := Username("Anonymous")
	if u != nil {
		username = u.Name
	}
	ret := GetAllUsersRegistered().PublicInfos()
	log.Printf("User %v queries all accounts.\n", username)
	return ret
}

// CancelAccount cancels(deletes) the User's own account
func (u *User) CancelAccount() error {
	if u == nil {
		return homererror.ErrNilUser
	}
	log.Printf("User %v canceled account.\n", u.Name)
	return nil
}

// SponsorMeeting creates a meeting, where User as the actor
func (u *User) SponsorMeeting(info MeetingInfo) (*Meeting, error) {
	if u == nil {
		return nil, homererror.ErrNilUser
	}
	m := NewMeeting(info)
	err := GetAllMeetings().Add(m)
	log.Printf("User %v sponsors meeting %v.\n", u.Name, info)
	return m, err
}

// AddParticipatorToMeeting just as its name
func (u *User) AddParticipatorToMeeting(meeting *Meeting, user *User) error {
	if u == nil {
		return homererror.ErrNilUser
	}

	err := meeting.Involve(user)
	log.Printf("User %v adds participator %v into Meeting %v.\n", u.Name, user.Name, meeting.Title)
	return err
}

// RemoveParticipatorFromMeeting just as its name
func (u *User) RemoveParticipatorFromMeeting(meeting *Meeting, user *User) error {
	if u == nil {
		return homererror.ErrNilUser
	}
	err := meeting.Exclude(user)
	log.Printf("User %v removes participator %v from Meeting %v.\n", u.Name, user.Name, meeting.Title)
	return err
}

// LogOut log out User's own (current working) account
func (u *User) LogOut() error {
	if u == nil {
		return homererror.ErrNilUser
	}

	log.Printf("User %v logs out.\n", u.Name)
	return nil
}

func (u *User) QueryMeetingByInterval(start, end time.Time) MeetingInfoListPrintable {
	// NOTE: FIXME: whatever, temporarily ignore the problem that the actor of query is Nil
	username := Username("Anonymous")
	if u != nil {
		username = u.Name
	}
	ret := u.involvedMeetings().Textualize()
	log.Printf("User %v queries meetings in time interval %v ~ %v.\n", username, start, end)
	return ret
}

func (u *User) meetingsSponsored() ([]*Meeting, error) {

	return nil, homererror.ErrNeedImplement
}

// CancelMeeting cancels(deletes) the given meeting which sponsored by User self, where User as the actor
func (u *User) CancelMeeting(meeting *Meeting) error {
	if u == nil {
		return homererror.ErrNilUser
	}

	err := meeting.Dissolve()
	log.Printf("User %v cancels Meeting %v.\n", u.Name, meeting.Title)
	return err
}

// QuitMeeting let User quit the given meeting, where User as the actor
func (u *User) QuitMeeting(meeting *Meeting) error {
	if u == nil {
		return homererror.ErrNilUser
	}

	err := meeting.Exclude(u)
	log.Printf("User %v quits Meeting %v.\n", u.Name, meeting.Title)
	return err
}
