package model

import (
	"config"

	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3" // dirver
	// _ "github.com/jinzhu/gorm/dialects/sqlite"  // dirver by gorm

	"util"
)

var userDB *gorm.DB

func init() {
	// db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")  // from gorm's doc
	db, err := gorm.Open("sqlite3", config.WorkingDir()+"agenda.db")
	util.PanicIf(err)
	userDB = db
}
