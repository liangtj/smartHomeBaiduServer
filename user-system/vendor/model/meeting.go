package model

import (
	"entity"
	"time"

	"util"
	log "util/logger"
)

func init() {
	m := &entity.MeetingInfoForDatabase{}
	if !userDB.HasTable(m) {
		err := userDB.CreateTable(m).Error
		util.PanicIf(err)

		// FIXME: wanted to model many2many relation, however, only without so the codes could work .....
		// u := &entity.UserInfoSerializable{}
		// err = agendaDB.Model(m).Related(u, "participations").Error
		util.PanicIf(err)

		log.Infof("\n ...... CreateTable %T. \n", m)
	}
}

// MeetingInfoAtomicService .
type MeetingInfoAtomicService struct{}

// MeetingInfoService .
var MeetingInfoService = MeetingInfoAtomicService{}

// Create .
func (*MeetingInfoAtomicService) Create(m *entity.MeetingInfoForDatabase) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Create(m).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Save .
func (*MeetingInfoAtomicService) Save(m *entity.MeetingInfoForDatabase) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Save(m).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Delete .
func (*MeetingInfoAtomicService) Delete(m *entity.MeetingInfoForDatabase) error {
	tx := userDB.Begin()
	util.PanicIf(tx.Error)

	if err := tx.Delete(m).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// FindAll .
func (*MeetingInfoAtomicService) FindAll() ([]entity.MeetingInfoForDatabase, error) {
	var rows []entity.MeetingInfoForDatabase
	err := userDB.Find(&rows).Error // CHECK: should check .Error ?
	return rows, err
}

// FindByTitle .
func (*MeetingInfoAtomicService) FindByTitle(title entity.MeetingTitle) (entity.MeetingInfoForDatabase, error) {
	var mInfo entity.MeetingInfoForDatabase
	err := userDB.First(&mInfo, entity.MeetingInfoForDatabase{Title: title}).Error
	return mInfo, err
}

// FindByInterval .
func (*MeetingInfoAtomicService) FindByInterval(start time.Time, end time.Time) ([]entity.MeetingInfoForDatabase, error) {
	var mInfos []entity.MeetingInfoForDatabase
	err := userDB.Where("start_time >= ? AND end_time <= ?", start, end).Find(&mInfos).Error
	return mInfos, err
}
