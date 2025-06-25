package db

import (
	"fmt"
	"time"
)

type BirthdayMsg struct {
	DiscordId     int
	BdayResponse  string
	TextChannelID string
}

type User struct {
	ID            int `gorm:"primaryKey"`
	DiscordId     int `gorm:"column:disc_id"`
	FirstName     string
	LastName      string
	Dob           time.Time
	BdayResponse  string
	TextChannelID string
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

func GetBirthdays(dateToCheck time.Time) ([]BirthdayMsg, error) {
	var users []User

	res := db.Where("EXTRACT(MONTH FROM dob) = ? AND EXTRACT(DAY FROM dob) = ?", int(dateToCheck.Month()), dateToCheck.Day()).Find(&users)
	if res.Error != nil {
		fmt.Println("Error: ", res.Error)
		return nil, res.Error
	}
	var bdays []BirthdayMsg
	for _, user := range users {
		bdays = append(bdays, BirthdayMsg{
			DiscordId:     user.DiscordId,
			BdayResponse:  user.BdayResponse,
			TextChannelID: user.TextChannelID,
		})
	}

	return bdays, nil
}
