package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/braydencw1/venova/db"
	"github.com/bwmarrin/discordgo"
)

type StreamInfo struct {
	Config     Config  `json:"config"`
	StartedAt  int64   `json:"started_at"`
	StreamCode string  `json:"stream_code"`
	Viewers    Viewers `json:"viewers"`
}

type Config struct {
	ClipBufferSecs int    `json:"clip_buffer_secs"`
	Clipping       bool   `json:"clipping"`
	VideoCodec     string `json:"video_codec"`
}

type Viewers struct {
	Total        int `json:"total"`
	WebRTC       int `json:"webrtc"`
	WebTransport int `json:"webtransport"`
}

type liveState struct {
	IsLive        bool
	LastStartedAt int64
	LastAlertTime time.Time
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func PollStreamer(s *discordgo.Session, list []db.HardlightStreamer) {
	state := make(map[int]*liveState)
	cooldown := time.Hour

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for _, streamer := range list {
			stream, live := checkLive(streamer.APIURL)

			st, exists := state[streamer.ID]
			if !exists {
				st = &liveState{}
				state[streamer.ID] = st
			}

			if !live {
				if st.IsLive {
					st.IsLive = false
					log.Printf("%s went offline", streamer.Name)
				}
				continue
			}

			newSession := st.LastStartedAt != 0 && st.LastStartedAt != stream.StartedAt
			if !st.IsLive || newSession {
				// cooldown protection
				if time.Since(st.LastAlertTime) < cooldown {
					continue
				}

				st.IsLive = true
				st.LastStartedAt = stream.StartedAt
				st.LastAlertTime = time.Now()

				handleLive(s, streamer, stream)
			}
		}
	}
}

func handleLive(s *discordgo.Session, streamer db.HardlightStreamer, stream StreamInfo) {
	started := time.Unix(stream.StartedAt, 0).Local()

	msg := fmt.Sprintf(
		"🔴 %s is LIVE!\nStarted: %s\nViewers: %d\n%s",
		streamer.Name,
		started.Format(time.Kitchen),
		stream.Viewers.Total,
		streamer.WatchURL,
	)

	if _, err := s.ChannelMessageSend(streamer.ChannelID, msg); err != nil {
		log.Print(err)
	}
}

func checkLive(url string) (StreamInfo, bool) {
	var stream StreamInfo

	resp, err := httpClient.Get(url)
	if err != nil {
		return stream, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return stream, false
	}

	if err := json.NewDecoder(resp.Body).Decode(&stream); err != nil {
		return stream, false
	}

	return stream, true
}
