package bot

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var AudioBuffer = make(chan []byte, 100)

func playAudioCmd(ctx CommandCtx) error {
	if ctx.Message.Author.ID == morthisId {
		gId := ctx.Message.GuildID
		vc := getUserVoiceChannel(ctx.Session, gId, ctx.Message.Author.ID)
		if vc == "" {
			return ctx.Reply("could not find vc")
		}
		voiceConn, err := ctx.Session.ChannelVoiceJoin(gId, vc, false, true)
		if err != nil {
			log.Fatalf("error joining vc to play audio: %s", err)
		}
		go monitorVoiceActivity(ctx.Session, gId, vc)
		go StartAudioReceiver(voiceConn)
		return ctx.Reply("Venova has enterd the chat...")
	}
	return nil
}

func StartAudioReceiver(vc *discordgo.VoiceConnection) {
	port := gatherVars()
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Printf("error resolving audio receiver addr: %s", err)
	}
	con, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("error listening to audio receiver: %s", err)
	}
	defer con.Close()

	buffer := make([]byte, 1024)

	for {
		n, _, err := con.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading connection stream UDP: %s", err)
			continue
		}
		audioData := buffer[:n]
		vc.Speaking(true)
		vc.OpusSend <- bytes.Clone(audioData)
		vc.Speaking(false)

	}

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
