package db

import (
	"log"
	"time"
)

type DndPlayDate struct {
	DateOfPlay        time.Time
	DndCampaignTcId   int64
	DndCampaignId     int
	DndCampaignRoleId int64
}

type DndCampaigns struct {
	Id     int
	Name   string
	TcId   int64
	RoleId int64
}

func GetPlayDates(dateToCheck time.Time) (bool, int64, int64, error) {
	var playdates DndPlayDate
	startOfDay := dateToCheck.Format("2006-01-02") + " 00:00:00"
	endOfDay := dateToCheck.Format("2006-01-02") + " 23:59:59"
	res := db.Where("date_of_play BETWEEN ? AND ?", startOfDay, endOfDay).Find(&playdates)
	if res.Error != nil {
		log.Printf("Error: %v", res.Error)
		return false, 0, 0, res.Error
	}
	if res.RowsAffected == 1 {
		return true, playdates.DndCampaignTcId, playdates.DndCampaignRoleId, nil
	} else {
		return false, 0, 0, nil
	}
}

func InsertPlayDate(playTime time.Time, roleId int64) (bool, error) {
	var playDateInfo DndPlayDate

	res := db.Where("dnd_campaign_role_id = ?", roleId).Order("dnd_campaign_role_id DESC").Limit(1).Find(&playDateInfo)

	if res.Error != nil {
		return false, res.Error
	}
	insertPlayDate := DndPlayDate{
		DateOfPlay:        playTime,
		DndCampaignTcId:   playDateInfo.DndCampaignTcId,
		DndCampaignId:     playDateInfo.DndCampaignId,
		DndCampaignRoleId: roleId,
	}
	insertRes := db.Create(&insertPlayDate)
	if insertRes.Error != nil {
		return false, res.Error
	}

	if res.RowsAffected == 1 {
		return true, nil
	} else {

		return false, nil
	}
}

func GetLatestPlayDate(dndRoleId int64) (time.Time, int64, error) {
	var playDate DndPlayDate
	res := db.Where("dnd_campaign_role_id = ?", dndRoleId).Order("date_of_play desc").First(&playDate)
	if res.Error != nil {
		return time.Time{}, 0, res.Error
	}
	return playDate.DateOfPlay, playDate.DndCampaignTcId, nil
}

func GetDndRoles() ([]int64, error) {
	var dndRole []DndCampaigns
	var roleArray []int64
	res := db.Find(&dndRole)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, campaign := range dndRole {
		roleArray = append(roleArray, campaign.RoleId)
	}
	return roleArray, nil
}
