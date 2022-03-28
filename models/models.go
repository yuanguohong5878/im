package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var db *gorm.DB

// Setup initializes the database instance
func Setup() {

	db, _ = gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/im_db?charset=utf8mb4&parseTime=True")

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName + "_tab"
	}
	db.SingularTable(true)
	db.LogMode(true)
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(64)
	sqlDB.SetMaxOpenConns(64)
	sqlDB.SetConnMaxLifetime(time.Minute)

}
