package model

import (
	"entity"

	"util"
	log "util/logger"
)

func init() {
	s := &entity.SessionInfo{}
	if !userDB.HasTable(s) {
		err := userDB.CreateTable(s).Error
		util.PanicIf(err)
		log.Infof("\n ...... CreateTable %T. \n", s)
	}
}

// SessionInfoAtomicService .
type SessionInfoAtomicService struct{}

// SessionInfoService .
var SessionInfoService = SessionInfoAtomicService{}

// Create .
func (*SessionInfoAtomicService) Create(s *entity.SessionInfo) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Create(s).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Save .
func (*SessionInfoAtomicService) Save(s *entity.SessionInfo) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Save(s).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete .
func (*SessionInfoAtomicService) Delete(s *entity.SessionInfo) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Delete(s).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// FindAll .
func (*SessionInfoAtomicService) FindAll() ([]entity.SessionInfo, error) {
	var rows []entity.SessionInfo
	err := userDB.Find(&rows).Error // CHECK: should check .Error ?
	return rows, err
}

// FindByToken .
func (*SessionInfoAtomicService) FindByToken(token entity.Token) (entity.SessionInfo, error) {
	var sInfo entity.SessionInfo
	err := userDB.First(&sInfo, entity.SessionInfo{Token: token}).Error
	return sInfo, err
}
