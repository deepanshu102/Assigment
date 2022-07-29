package database

import (
	"log"

	"github.com/deep/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	migration()
}

const dns = "root:root@tcp(127.0.0.1:3306)/dummy?charset=utf8mb4&parseTime=True&loc=Local"

func ConnectionStablish() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Connection Error", err)
	}

	return db, err
}
func migration() {
	db, _ := ConnectionStablish()
	db.AutoMigrate(&models.Users{})
}
