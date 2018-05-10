package entity

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"
)

type Token string

// Empty checks if Token empty
func (token Token) Empty() bool {
	return token == ""
}

// Valid checks if Token valid
func (token Token) Valid() bool {
	return !token.Empty() // NOTE: may not only !empty
}

func (token Token) String() string {
	return string(token)
}

var base = 0

func randomBytesGen(n int) ([]byte, error) {
	ret := make([]byte, n)
	_, err := rand.Read(ret)
	return ret, err
}

func TokenGen(n int) Token {
	bytes, err := randomBytesGen(n)
	i, retryMaxCount := 0, 100
	for err != nil && i < retryMaxCount {
		bytes, err = randomBytesGen(n)
	}
	if i == retryMaxCount {
		log.Fatalf("Fail to generate random-bytes, error: %q\n", err.Error())
	}

	s := base64.URLEncoding.EncodeToString(bytes)
	return Token(s)
}

type SessionInfo struct {
	Token     Token                `gorm:"primary_key"`
	ExpiredAt time.Time            `gorm:"not NULL;column:expired_at"`
	User      UserInfoSerializable `gorm:"not NULL"`
}

type Session struct {
	SessionInfo
}

func (sess *Session) Valid() bool {
	return time.Now().Before(sess.ExpiredAt)
}

func (sess *Session) Destroy() {
	sess.Token = ""
}

func (sess *Session) Reset(newTime time.Time) {
	sess.ExpiredAt = newTime
	if !sess.Valid() {
		sess.Destroy()
	}
}
