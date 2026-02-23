package db

import (
	"time"
)

type Streamer struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	URL         string `gorm:"uniqueIndex"`
	IsLive      bool
	LastChecked time.Time
}

func GetStreamers() ([]Streamer, error) {
	var streamers []Streamer

	err := db.Find(&streamers).Error
	if err != nil {
		return []Streamer{}, err
	}
	return streamers, nil
}
