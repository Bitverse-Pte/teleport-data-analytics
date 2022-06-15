package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(url string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(url))
	if err != nil {
		logrus.Fatalln("db initing fail", err)
		return nil
	}
	return db
}
