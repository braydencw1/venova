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
	Id        int `gorm:"primaryKey"`
	DiscordId int `gorm:"column:disc_id"`
	FirstName string
	LastName  string
	Dob       time.Time
}

func GetBirthdays(dateToCheck time.Time) ([]int, error) {
	var users []User

	res := db.Where("EXTRACT(MONTH FROM dob) = ? AND EXTRACT(DAY FROM dob) = ?", int(dateToCheck.Month()), dateToCheck.Day()).Find(&users)
	if res.Error != nil {
		fmt.Println("Error: ", res.Error)
		return nil, res.Error
	}
	var discIds []int
	for _, user := range users {
		log.Printf("Today's birthdays are: %v", user)
		discIds = append(discIds, user.DiscordId)
	}
	return discIds, nil
}
