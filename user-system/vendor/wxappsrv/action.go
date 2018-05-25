package wxappsrv

import (
	errors "convention/errors"
	"entity"
	"model"
	log "util/logger"

	"github.com/jinzhu/gorm"
)

func QueryAccountByUsername(name entity.UserIdentifier) (entity.UserInfo, error) {
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
