package db

import (
	"fmt"
	"log"
	"time"
)

type DndPlayDate struct {
	DateOfPlay        time.Time
	DndCampaignTcId   int64
	DndCampaignId     int
	DndCampaignRoleId string
}

type DndCampaigns struct {
	Id     int
	Name   string
	TcId   int64
	RoleId string
}

func GetPlayDates(dateToCheck time.Time) (bool, int64, string, error) {
	var playdates DndPlayDate
	res := db.Where("DATE_TRUNC('day', date_of_play::date) = DATE_TRUNC('day', ?::date)", dateToCheck).Take(&playdates)
	if res.Error != nil {
		log.Printf("Error: %v", res.Error)
		return false, 0, "", res.Error
	}
	return true, playdates.DndCampaignTcId, playdates.DndCampaignRoleId, nil
}

func InsertPlayDate(playTime time.Time, roleId string) error {
	var playDateInfo DndPlayDate

	res := db.Where("dnd_campaign_role_id = ?", roleId).Take(&playDateInfo)
	if res.Error != nil {
		return res.Error
	}

	insertPlayDate := DndPlayDate{
		DateOfPlay:        playTime,
		DndCampaignTcId:   playDateInfo.DndCampaignTcId,
		DndCampaignId:     playDateInfo.DndCampaignId,
		DndCampaignRoleId: roleId,
	}
	insertRes := db.Create(&insertPlayDate)
	if insertRes.Error != nil {
		return res.Error
	}
	return nil
}

func GetLatestPlayDate(dndRoleId string) (time.Time, int64, error) {
	var playDate DndPlayDate
	res := db.Where("dnd_campaign_role_id = ?", dndRoleId).Order("date_of_play desc").First(&playDate)
	if res.Error != nil {
		return time.Time{}, 0, res.Error
	}
	return playDate.DateOfPlay, playDate.DndCampaignTcId, nil
}

func GetDndRoles() ([]string, error) {
	var roleArray []string
	var dndRole []DndCampaigns
	res := db.Find(&dndRole)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, campaign := range dndRole {
		roleArray = append(roleArray, campaign.RoleId)
	}
	return roleArray, nil
}

type DndResponses struct {
	Id       int `gorm:"primaryKey"`
	Response string
}

func DndMsgResponse() string {
	var query DndResponses
	res := db.Order("RANDOM()").Find(&query)
	if res.Error != nil {
		err := fmt.Sprintf("Error, %s", res.Error)
		return err
	}
	return query.Response

}
