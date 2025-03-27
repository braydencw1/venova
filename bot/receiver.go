package bot

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type AudioReceiver struct {
	StopChan chan bool
	Wg       sync.WaitGroup
	Vc       *discordgo.VoiceConnection
	Conn     *net.UDPConn
}

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
		ar := NewAudioReceiver(voiceConn)
		activeReceiver = ar
		go monitorVoiceActivity(s, gId, vc)
		go ar.Start()
		return ctx.Reply("Venova has entered the chat...")
	}
	return nil
}

func NewAudioReceiver(vc *discordgo.VoiceConnection) *AudioReceiver {
	return &AudioReceiver{
		StopChan: make(chan bool, 1),
		Vc:       vc,
	}
}

func gatherVars() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	return GetEnvOrDefault("AUDIO_SERVER_PORT", "5005")
}

func (ar *AudioReceiver) Start() {
	port := gatherVars()
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}
	ar.Conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error listening on UDP: %v", err)
	}
	buffer := make([]byte, 4000)
	err = ar.Vc.Speaking(true)
	if err != nil {
		log.Printf("Error to speak: %s", err)
	}
	ar.Wg.Add(1)
	defer func() {
		ar.Cleanup()
	}()

	for {
		select {
		case <-ar.StopChan:
			return
		default:
			n, _, err := ar.Conn.ReadFrom(buffer)
			if err != nil {
				log.Printf("Error reading UDP: %v", err)
				continue
			}
			if n <= 12 {
				continue
			}

			opusFrame := make([]byte, n-12)
			copy(opusFrame, buffer[12:n])
			ar.Vc.OpusSend <- opusFrame
		}

	}
}

func (ar *AudioReceiver) Cleanup() {
	err := ar.Vc.Speaking(false)
	if err != nil {
		log.Printf("Unable to defer vc speak false: %s", err)
	}
	if err := ar.Vc.Disconnect(); err != nil {
		log.Printf("Unable to close vc speak: %s", err)
	}
	if ar.Conn != nil {
		log.Printf("Closing")
		ar.Conn.Close()
		ar.Conn = nil
	}
	activeReceiver = nil
	log.Printf("Audio Receiver Stopped.")
	ar.Wg.Done()
}

func (ar *AudioReceiver) Stop() {
	select {
	case ar.StopChan <- true:
	default:
		log.Printf("Stopping of Audio Receiver already in progress.")
	}
	ar.Wg.Wait()
}
