package db

import (
	"log"
	"time"
)

type DndPlayDate struct {
	DateOfPlay time.Time
}

func GetPlayDates(dateToCheck time.Time) (bool, error) {
	var playdates DndPlayDate
	startOfDay := dateToCheck.Format("2006-01-02") + " 00:00:00"
	endOfDay := dateToCheck.Format("2006-01-02") + " 23:59:59"
	res := db.Where("date_of_play BETWEEN ? AND ?", startOfDay, endOfDay).Find(&playdates)
	if res.Error != nil {
		log.Printf("Error: %v", res.Error)
		return false, res.Error
	}
	if res.RowsAffected == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func InsertPlayDate(playTime time.Time) (bool, error) {
	playDate := DndPlayDate{
		DateOfPlay: playTime,
	}
	res := db.Create(&playDate)
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
func GetLatestPlayDate() (time.Time, error) {
	var playDate DndPlayDate
	res := db.Last(&playDate)
	if res.Error != nil {
		return time.Time{}, res.Error
	}
	return playDate.DateOfPlay, nil
}
