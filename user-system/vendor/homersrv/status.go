package homersrv

import (
	"config"
	errors "convention/homererror"
	"io/ioutil"
	"os"
	log "util/logger"
)

var loginedUser = Username("")

func LoginedUser() *User {
	name := loginedUser
	return name.RefInAllUsers()
}

func LoadLoginStatus() {
	buf, err := ioutil.ReadFile(config.UserLoginStatusPath())
	// fin, err := os.Open(config.UserLoginStatusPath())
	// defer fin.Close()
	if err != nil {
		log.Fatal(err)
	}
	loginedUser = Username(string(buf))
	// fin.Read()
	// return errors.ErrNeedImplement
}

func SaveLoginStatus() error {
	fout, err := os.Create(config.UserLoginStatusPath())
	defer fout.Close()
	if err != nil {
		log.Fatal(err)
	}
	fout.WriteString(string(loginedUser))
	return errors.ErrNeedImplement
}
