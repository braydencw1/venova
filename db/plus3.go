package db

type Streamer struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	URL       string `gorm:"uniqueIndex"`
	ChannelID string
	IsLive    bool
}

func GetStreamers() ([]Streamer, error) {
	var streamers []Streamer

	err := db.Find(&streamers).Error
	if err != nil {
		return []Streamer{}, err
	}
	return streamers, nil
}
