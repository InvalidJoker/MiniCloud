package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database = gorm.DB

func NewDatabase() (*Database, error) {
	db, err := gorm.Open(sqlite.Open("data/config/data.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Server{})

	return db, nil
}
