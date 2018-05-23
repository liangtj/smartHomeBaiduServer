package homersrv

import (
	errors "convention/homererror"
	"entity"
	"model"
	log "util/logger"

	"github.com/jinzhu/gorm"
)

// RegisterUser ...
func RegisterUser(uInfo entity.UserInfo) error {
	if !uInfo.Name.Valid() {
		return errors.ErrInvalidUsername
	}

	u := entity.NewUser(uInfo)
	if err := model.UserInfoService.Create(&uInfo); err != nil {
		log.Error(err) // TODO: should not be like this
	}
	err := entity.GetAllUsersRegistered().Add(u)
	return err
}

func LogIn(name Username, auth Auth) error {
	u := name.RefInAllUsers()
	if u == nil {
		return errors.ErrNilUser
	}

	log.Printf("User %v logs in.\n", name)

	if LoginedUser() != nil {
		return errors.ErrLoginedUserAuthority
	}

	if verified := u.Auth.Verify(auth); !verified {
		return errors.ErrFailedAuth
	}

	loginedUser = name

	return nil
}

// LogOut log out User's own (current working) account
// TODO:
func LogOut(name Username) error {
	u := name.RefInAllUsers()

	// check if under login status, TODO: check the login status
	if logined := LoginedUser(); logined == nil {
		return errors.ErrUserNotLogined
	} else if logined != u {
		return errors.ErrUserAuthority
	}

	err := u.LogOut()
	if err != nil {
		log.Errorf("Failed to log out, error: %q.\n", err.Error())
	}
	loginedUser = ""
	return err
}

// QueryAccountAll queries all accounts
func QueryAccountAll() []UserInfoPublic {
	// NOTE: FIXME: whatever, temporarily ignore the problem that the actor of query is Nil
	// Hence, now if so, agenda would crash for `Nil.Name`
	ret := LoginedUser().QueryAccountAll()
	return ret
}

// CancelAccount cancels(deletes) LoginedUser's account
func CancelAccount() error {
	u := LoginedUser()
	if u == nil {
		return errors.ErrUserNotLogined
	}

	if err := entity.GetAllMeetings().ForEach(func(m *Meeting) error {
		if m.SponsoredBy(u.Name) {
			return m.Dissolve()
		}
		if m.ContainsParticipator(u.Name) {
			return m.Exclude(u)
		}
		return nil
	}); err != nil {
		log.Error(err)
	}

	if err := entity.GetAllUsersRegistered().Remove(u); err != nil {
		log.Error(err)
	}
	if err := u.LogOut(); err != nil {
		log.Error(err)
	}

	err := u.CancelAccount()
	return err
}

// ----------------------------------------------------------------
// @@binly: new

func QueryAccountByUsername(name entity.Username) (entity.UserInfo, error) {
	if !name.Valid() {
		return entity.UserInfo{}, errors.ErrInvalidUsername
	}
	uInfo, err := model.UserInfoService.FindByUsername(name)
	return uInfo, err
}
func Authorize(token entity.Token) (entity.SessionInfo, error) {
	if !token.Valid() {
		return entity.SessionInfo{}, ErrInvalidToken
	}
	sInfo, err := model.SessionInfoService.FindByToken(token)
	if err != nil {
		return sInfo, err
	}

	sess := entity.Session{sInfo}
	if !sess.Valid() {
		if err := model.SessionInfoService.Delete(&sInfo); err != nil {
			log.Errorf("Delete a bad session failed, error: %q\n", err.Error())
		}
		return entity.SessionInfo{}, ErrSessionExpired
	}

	return sInfo, err
}
func DeleteSession(sInfo *entity.SessionInfo) error {

	if err := model.SessionInfoService.Delete(sInfo); err != nil {
		log.Errorf("Failed to delete session(Token:%q), error: %q.\n", sInfo.Token, err.Error())
		return err

		// if err = sInfo.User.LogOut(); err != nil {
		// 	log.Errorf("Failed to log out, error: %q.\n", err.Error())
		// 	return err
		// }
	}
	return nil
}

// TODO: FIXME: limit the number of sessions a User can create
func CreateSession(sInfo *entity.SessionInfo) error {
	token := entity.TokenGen(32)
	_, err := model.SessionInfoService.FindByToken(token)
	i, retryMaxCount := 0, 100
	for err != gorm.ErrRecordNotFound && i < retryMaxCount {
		token = entity.TokenGen(32)
		_, err = model.SessionInfoService.FindByToken(token)
	}
	if i == retryMaxCount {
		log.Fatalf("Fail to generate a new token, error: %q\n", err.Error())
	}

	sInfo.Token = token
	err = model.SessionInfoService.Create(sInfo)
	return err
}
