package bot

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
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

		//go monitorVoiceActivity(s, gId, vc)
		go StartAudioReceiver(voiceConn)
		return nil
		// return ctx.Reply("Venova has entered the chat...")
	}
	return nil
}

// func StartAudioReceiver(vc *discordgo.VoiceConnection) {
// 	port := gatherVars()
// 	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%s", port))
// 	if err != nil {
// 		log.Printf("error resolving audio receiver addr: %s", err)
// 	}
// 	log.Printf("port here:%s", port)
// 	log.Printf("addr here:%s", addr)
// 	con, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		log.Printf("error listening to audio receiver: %s", err)
// 		return
// 	}
// 	defer con.Close()

// 	buffer := make([]byte, 1200)

// 	for {
// 		n, _, err := con.ReadFrom(buffer)
// 		if err != nil {
// 			log.Printf("Error reading connection stream UDP: %s", err)
// 			continue
// 		}
// 		log.Printf("Received %d bytes", n)
// 		audioData := buffer[:n]
// 		vc.Speaking(true)
// 		vc.OpusSend <- bytes.Clone(audioData)
// 		vc.Speaking(false)

// 	}

// }

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
	defer conn.Close()

	buffer := make([]byte, 4000)
	vc.Speaking(true)
	defer vc.Speaking(false)

	for {
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

//ffmpeg -f dshow -i audio="CABLE Output (VB-Audio Virtual Cable)" -ac 2 -ar 48000 -c:a libopus -b:a 64k -frame_size 960 -f rtp rtp://<ip>:<port>
