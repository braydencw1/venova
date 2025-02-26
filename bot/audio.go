package bot

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	stopAudioReceiver = make(chan bool, 1)
	wg                sync.WaitGroup
)

func playAudioCmd(ctx CommandCtx) error {
	if ctx.Message.Author.ID == morthisId {
		gId := ctx.Message.GuildID
		s := ctx.Session
		vc := getUserVoiceChannel(s, gId, ctx.Message.Author.ID)
		if vc == "" {
			return ctx.Reply("Could not find a voice channel.")
		}

		voiceConn, err := s.ChannelVoiceJoin(gId, vc, false, true)
		if err != nil {
			log.Printf("Error joining VC to play audio: %s", err)
			return ctx.Reply("Failed to join voice channel.")
		}

		go monitorVoiceActivity(s, gId, vc)
		go StartAudioReceiver(voiceConn)
		return ctx.Reply("Venova has entered the chat...")
	}
	return nil
}

func gatherVars() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	port := os.Getenv("AUDIO_SERVER_PORT")
	if port == "" {
		log.Printf("AUDIO_SERVER_PORT not defined, defaulting to 5005")
		port = "5005"
	}
	return port
}

func StartAudioReceiver(vc *discordgo.VoiceConnection) {
	port := gatherVars()
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error listening on UDP: %v", err)
	}

	buffer := make([]byte, 4000)
	err = vc.Speaking(true)
	if err != nil {
		log.Printf("Error to speak: %s", err)
	}
	wg.Add(1)
	defer func() {
		err := vc.Speaking(false)
		if err != nil {
			log.Printf("Unable to defer vc speak false: %s", err)
		}
		err = vc.Disconnect()
		if err != nil {
			log.Printf("Unable to close vc speak: %s", err)
		}

		conn.Close()
		log.Printf("Audio Receiver Stopped.")
		wg.Done()
	}()

	for {
		select {
		case <-stopAudioReceiver:
			return
		default:
			n, _, err := conn.ReadFrom(buffer)
			if err != nil {
				log.Printf("Error reading UDP: %v", err)
				continue
			}
			// Ensure packet is large enough to contain RTP header
			if n <= 12 {
				continue
			}

			// Strip the RTP header (first 12 bytes)
			opusFrame := make([]byte, n-12)
			copy(opusFrame, buffer[12:n])
			vc.OpusSend <- opusFrame
		}

	}
}

func stopAudio() {
	stopAudioReceiver <- true
}
