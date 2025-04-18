package db

import (
	"fmt"
	"log"
	"time"
)

type User struct {
	ID           int `gorm:"primaryKey"`
	DiscordId    int `gorm:"column:disc_id"`
	FirstName    string
	LastName     string
	Dob          time.Time
	BdayResponse string
}

type AdminUser struct {
	UserID int  `gorm:"primaryKey"`
	User   User `gorm:"foreignKey:UserID;references:ID"`
}

type McAdminUser struct {
	UserID int  `gorm:"primaryKey"`
	User   User `gorm:"foreignKey:UserID;references:ID"`
}

type BirthdayReminderUser struct {
	UserID int  `gorm:"primaryKey"`
	User   User `gorm:"foreignKey:UserID;references:ID"`
}

func GetBirthdays(dateToCheck time.Time) (map[int]string, error) {
	var users []User

	res := db.Where("EXTRACT(MONTH FROM dob) = ? AND EXTRACT(DAY FROM dob) = ?", int(dateToCheck.Month()), dateToCheck.Day()).Find(&users)
	if res.Error != nil {
		fmt.Println("Error: ", res.Error)
		return nil, res.Error
	}
	bdayMap := make(map[int]string)

	for _, user := range users {
		log.Printf("Today's birthdays are: %v", user)
		bdayMap[user.DiscordId] = user.BdayResponse
	}
	return bdayMap, nil
}
