package model

import (
	"entity"

	"util"
	log "util/logger"
)

func init() {
	u := &entity.UserInfoSerializable{}
	if !userDB.HasTable(u) {
		err := userDB.CreateTable(u).Error
		util.PanicIf(err)
		log.Infof("\n ...... CreateTable %T. \n", u)
	}
}

// UserInfoAtomicService .
type UserInfoAtomicService struct{}

// UserInfoService .
var UserInfoService = UserInfoAtomicService{}

// func loadAllUserFromDB(db *DB) {}

// Create .
func (*UserInfoAtomicService) Create(u *entity.UserInfoSerializable) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Save .
func (*UserInfoAtomicService) Save(u *entity.UserInfoSerializable) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Save(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete .
func (*UserInfoAtomicService) Delete(u *entity.UserInfoSerializable) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Delete(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// FindAll .
func (*UserInfoAtomicService) FindAll() ([]entity.UserInfoSerializable, error) {
	var rows []entity.UserInfoSerializable
	err := userDB.Find(&rows).Error
	return rows, err
}

// FindByUsername .
func (*UserInfoAtomicService) FindByUsername(name entity.Username) (entity.UserInfoSerializable, error) {
	var uInfo entity.UserInfoSerializable

	// agendaDB.First(&uInfo, entity.UserInfoSerializable{Name: name}) // TODEL: sad to anonymous member ...
	u := entity.UserInfoSerializable{}
	u.ID = name
	err := userDB.First(&uInfo, u).Error
	return uInfo, err
}
