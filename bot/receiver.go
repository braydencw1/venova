package bot

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
)

func init() {
	if err := gotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

type AudioReceiver struct {
	done chan struct{}
	vc   *discordgo.VoiceConnection
	conn *net.UDPConn
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
		port := getAudioServerPort()
		if err := ensureValidPort(port); err != nil {
			log.Fatalf("%s", err)
		}
		ar, err := NewAudioReceiver(voiceConn, port)
		if err != nil {
			log.Fatalf("%s", err)
		}
		activeReceiver = ar
		go monitorVoiceActivity(s, gId, vc)

		go ar.Run()
		return ctx.Reply("Venova has entered the chat...")
	}
	return nil
}

func ensureValidPort(port string) error {
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return fmt.Errorf("invalid port: %s", port)
	}
	return nil
}

func NewAudioReceiver(vc *discordgo.VoiceConnection, port string) (*AudioReceiver, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, fmt.Errorf("Error resolving UDP address: %w", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("Error listening on UDP: %w", err)
	}
	if err := vc.Speaking(true); err != nil {
		return nil, fmt.Errorf("Error to speak: %w", err)
	}
	return &AudioReceiver{
		done: make(chan struct{}),
		vc:   vc,
		conn: conn,
	}, nil
}

func getAudioServerPort() string {
	return GetEnvOrDefault("AUDIO_SERVER_PORT", "5005")
}

func (ar *AudioReceiver) Run() {
	buffer := make([]byte, 4000)
	defer ar.Cleanup()
	for {
		select {
		case <-ar.done:
			return
		default:
			n, _, err := ar.conn.ReadFrom(buffer)
			if err != nil {
				log.Printf("Error reading UDP: %v", err)
				continue
			}
			if n <= 12 {
				continue
			}

			opusFrame := make([]byte, n-12)
			copy(opusFrame, buffer[12:n])
			ar.vc.OpusSend <- opusFrame
		}

	}
}

func (ar *AudioReceiver) Cleanup() {
	err := ar.vc.Speaking(false)
	if err != nil {
		log.Printf("Unable to defer vc speak false: %s", err)
	}
	if err := ar.vc.Disconnect(); err != nil {
		log.Printf("Unable to close vc speak: %s", err)
	}
	if ar.conn != nil {
		ar.conn.Close()
	}
	log.Println("Audio Receiver stopped.")
}

func (ar *AudioReceiver) Stop() {
	close(ar.done)
}
