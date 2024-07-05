package db

import (
	"fmt"
	"log"
	"time"
)

type User struct {
	Id           int `gorm:"primaryKey"`
	DiscordId    int `gorm:"column:disc_id"`
	FirstName    string
	LastName     string
	Dob          time.Time
	BdayResponse string
}

func GetBirthdays(dateToCheck time.Time) (map[int]string, error) {
	var users []User

	res := db.Where("EXTRACT(MONTH FROM dob) = ? AND EXTRACT(DAY FROM dob) = ?", int(dateToCheck.Month()), dateToCheck.Day()).Find(&users)
	if res.Error != nil {
		fmt.Println("Error: ", res.Error)
		return nil, res.Error
	}
	discio := make(map[int]string)

	for _, user := range users {
		log.Printf("Today's birthdays are: %v", user)
		discio[user.DiscordId] = user.BdayResponse
	}
	return discio, nil
}
