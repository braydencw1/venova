package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDatabase(dsn string) error {
	it, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db = it

	err = db.AutoMigrate(
		&User{},
	)
	if err != nil {
		return err
	}

	return nil
}
