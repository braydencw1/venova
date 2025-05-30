package bot

import (
	"fmt"

	"github.com/braydencw1/venova/db"
	"gorm.io/gorm"
)

type IdentityChecker interface {
	IsAdmin(uID string) bool
	IsMinecraftAdmin(uID string) bool
	WantsBirthdayReminder() ([]string, error)
}

type DBIdentityChecker struct {
	DB *gorm.DB
}

var checker IdentityChecker

func GetIdentityChecker() IdentityChecker {
	checker = &DBIdentityChecker{DB: db.GetDB()}
	return checker
}

func (c *DBIdentityChecker) WantsBirthdayReminder() ([]string, error) {
	var remindUsers []db.BirthdayReminderUser
	err := c.DB.Preload("User").Find(&remindUsers).Error
	if err != nil {
		return nil, err
	}
	discordIDs := []string{}
	for _, r := range remindUsers {
		discordIDs = append(discordIDs, fmt.Sprint(r.User.DiscordId))
	}
	return discordIDs, nil
}

// Fix this to Merge with Minecraft Admin
func (c *DBIdentityChecker) IsAdmin(uID string) bool {
	var user db.User
	if err := c.DB.First(&user, "disc_id =?", uID).Error; err != nil {
		return false
	}

	var adminUser db.AdminUser
	if err := c.DB.First(&adminUser, "user_id = ?", user.ID).Error; err != nil {
		return false
	}

	return true
}

// See above comment nerd
func (c *DBIdentityChecker) IsMinecraftAdmin(uID string) bool {
	var user db.User
	if err := c.DB.First(&user, "disc_id =?", uID).Error; err != nil {
		return false
	}

	var mcAdminUser db.McAdminUser
	if err := c.DB.First(&mcAdminUser, "user_id = ?", user.ID).Error; err != nil {
		return false
	}

	return true
}
