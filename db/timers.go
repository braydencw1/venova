package db

import "time"

type Timer struct {
	ID          int    `gorm:"primaryKey"`
	TargetID    string // discord user ID to DM when the timer fires
	CreatedByID string // discord user ID that set the timer
	FiresAt     time.Time
}

func InsertTimer(t *Timer) error {
	return db.Create(t).Error
}

func GetPendingTimers() ([]Timer, error) {
	var timers []Timer
	res := db.Find(&timers)
	return timers, res.Error
}

func DeleteTimer(id int) error {
	return db.Delete(&Timer{}, id).Error
}
