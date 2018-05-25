package entity

import (
	"convention/codec"
)

var allMeetings *MeetingList

func (title MeetingTitle) RefInAllMeetings() *Meeting {
	return allMeetings.Ref(title)
}
func GetAllMeetings() *MeetingList {
	return allMeetings
}

// LoadAllMeeting concretely loads all Meetings
func LoadAllMeeting(decoder codec.Decoder) {
	allMeetings = LoadedMeetingList(decoder)
}

// SaveAllMeeting concretely saves all Meetings
func SaveAllMeeting(encoder codec.Encoder) error {
	meetings := GetAllMeetings()
	return meetings.Save(encoder)
}
