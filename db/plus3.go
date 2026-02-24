package db

type HardlightStreamer struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	APIURL    string `gorm:"uniqueIndex"`
	WatchURL  string
	ChannelID string
	IsLive    bool
}

func GetHardlightStreamer() ([]HardlightStreamer, error) {
	var streamer []HardlightStreamer

	err := db.Find(&streamer).Error
	if err != nil {
		return []HardlightStreamer{}, err
	}

	return streamer, nil
}
