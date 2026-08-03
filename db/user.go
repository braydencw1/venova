package db

import (
	"fmt"
	"time"
)

type BirthdayMsg struct {
	DiscordId     int64
	BdayResponse  string
	TextChannelID string
}

type User struct {
	ID        int   `gorm:"primaryKey"`
	DiscordId int64 `gorm:"column:disc_id"`
	FirstName string
	LastName  string
	// A birthday is a calendar day, not an instant. Storing it as timestamptz
	// made EXTRACT(MONTH/DAY) depend on the session timezone, which shifted
	// birthdays by a day. Keep this as date so no conversion happens.
	Dob           *time.Time `gorm:"type:date"`
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

	res := db.Where("dob IS NOT NULL AND EXTRACT(MONTH FROM dob) = ? AND EXTRACT(DAY FROM dob) = ?", int(dateToCheck.Month()), dateToCheck.Day()).Find(&users)
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
