package entity

import (
	"convention/codec"
)

// CHECK: Not sure where to place ...
var allUsersRegistered *UserList

// RefInAllUsers returns the ref of a Registered User depending on the Username
func (name UserIdentifier) RefInAllUsers() *User {
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
