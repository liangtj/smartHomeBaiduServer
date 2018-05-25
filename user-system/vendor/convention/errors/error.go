package errors

import "errors"

var New = errors.New

var (
	ErrNeedImplement = errors.New("this function need to be implemented")

	// User
	ErrNilUser           = errors.New("a nil user/*user is to be used")
	ErrExistedUser       = errors.New("the user has been existed")
	ErrUserNotFound      = errors.New("cannot find the user")
	ErrUserNotRegistered = errors.New("cannot find the user for registered one")
	ErrUserNotLogined    = errors.New("no user logined")

	ErrEmptyUsername   = errors.New("given username cannot be empty")
	ErrInvalidUsername = errors.New("given username is invalid")

	ErrFailedAuth           = errors.New("auth verify failed.")
	ErrLoginedUserAuthority = errors.New("only the User has logged in can modify login status.")

	ErrUserAuthority = errors.New("only the User self can modify his/her account")

	ErrNilSponsor            = errors.New("the sponsor cannot be nil")
	ErrSponsorAuthority      = errors.New("only the sponsor can modify the meeting")
	ErrSponsorResponsibility = errors.New("the sponsor can only cancel but not quit the meeting")

	// Meeting
	ErrNilMeeting          = errors.New("a nil meeting/*meeting is to be used")
	ErrExistedMeeting      = errors.New("the meeting has been existed")
	ErrExistedMeetingTitle = errors.New("the meeting title has been existed")
	ErrMeetingNotFound     = errors.New("cannot find the meeting")

	ErrEmptyMeetingTitle   = errors.New("given meeting title cannot be empty")
	ErrInvalidMeetingTitle = errors.New("given meeting title is invalid")

	// Time
	// InvalidTime         = errors.New("startTime/EndTime is not valid")
	ErrInvalidTimeInterval    = errors.New("the EndTime must be after StartTime")
	ErrConflictedTimeInterval = errors.New("given time interval conflicts with existed interval")

	// Information
	ErrGivenConflictedInfo = errors.New("given a not reasonable information")
)

type WxappError struct {
	msg string
}

func NewWxappError(msg string) *WxappError {
	return &WxappError{
		msg: msg,
	}
}

func (e *WxappError) Error() string {
	return e.msg
}
