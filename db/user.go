package db

import (
	"fmt"
	"log"
	"time"
)

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
