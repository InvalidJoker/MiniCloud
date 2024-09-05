package database

import (
	"fmt"
	"minicloud/internal/config"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database = gorm.DB

func NewDatabase(config2 config.Config) (*Database, error) {
	url := config2.DatabaseURL

	if url == "" {
		return nil, fmt.Errorf("database url is empty")
	}

	// check if url starts with mysql, mariadb, postgres

	var db *gorm.DB
	var err error

	if config2.DatabaseType == "mysql" || strings.HasPrefix(url, "mysql") || strings.HasPrefix(url, "mariadb") {
		db, err = NewMySQLDatabase(config2)
	} else if config2.DatabaseType == "postgres" || strings.HasPrefix(url, "postgres") {
		db, err = NewPostgresDatabase(config2)
	} else {
		return nil, fmt.Errorf("unknown database type")
	}

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Server{}, &Template{})

	return db, nil
}

func NewPostgresDatabase(config2 config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config2.DatabaseURL), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewMySQLDatabase(config2 config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config2.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
