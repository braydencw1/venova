package bot

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/braydencw1/venova/db"
	"github.com/bwmarrin/discordgo"
)

type StreamData struct {
	Config struct {
		ClipBufferSecs int    `json:"clip_buffer_secs"`
		Clipping       bool   `json:"clipping"`
		VideoCodec     string `json:"video_codec"`
	} `json:"config"`
	StartedAt  int64  `json:"started_at"`
	StreamCode string `json:"stream_code"`
	Viewers    struct {
		Total        int `json:"total"`
		WebRTC       int `json:"webrtc"`
		WebTransport int `json:"webtransport"`
	} `json:"viewers"`
}

func PollStreamer(s *discordgo.Session, list []db.Streamer) {

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for _, streamer := range list {
			live, err := checkLive(streamer.URL)
			if err != nil {
				log.Printf("error checking live report: %s", err)
			}

			if live {
				msg := fmt.Sprintf("🔴 %s is LIVE!", streamer.Name)
				_, err := s.ChannelMessageSend(streamer.ChannelID, msg)
				if err != nil {
					log.Print(err)
				}
			}

		}
	}
}

func checkLive(apiURL string) (bool, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return false, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	switch resp.StatusCode {

	case http.StatusOK:
		return true, nil

	case http.StatusNotFound:
		return false, nil

	default:
		return false, nil
	}
}
