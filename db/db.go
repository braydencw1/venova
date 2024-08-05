package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func OpenDatabase(dsn string) error {
	it, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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
