package homersrv

var loginedUser = Username("")

func LoginedUser() *User {
	name := loginedUser
	return name.RefInAllUsers()
}
