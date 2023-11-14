package db

import (
	"fmt"
	"log"
	"time"

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

type User struct {
	ID         int `gorm:"primaryKey"`
	Disc_id    int
	First_name string
	Last_name  string
	Dob        time.Time
}

func GetBirthdays(dateToCheck time.Time) ([]int, error) {
	if db == nil {
		return nil, fmt.Errorf("database error")
	}

	var users []User

	res := db.Where("EXTRACT(MONTH FROM dob) = ? AND EXTRACT(DAY FROM dob) = ?", int(dateToCheck.Month()), dateToCheck.Day()).Find(&users)
	if res.Error != nil {
		fmt.Println("Error: ", res.Error)
		return nil, res.Error
	}
	var discIds []int
	for _, user := range users {
		log.Printf("Today's birthdays are: %v", user)
		discIds = append(discIds, user.Disc_id)
	}
	return discIds, nil
}
