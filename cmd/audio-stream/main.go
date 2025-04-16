package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
	"github.com/braydencw1/venova"
	"github.com/braydencw1/venova/bot"
	"github.com/joho/godotenv"
)

var cli struct {
	Version bool `help:"Show version" short:"v"`
}

func main() {
	kong.Parse(&cli)
	if cli.Version {
		ver, err := venova.GetVersionInfo("venova-audio-stream")
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Println(ver)
		os.Exit(0)
	}
	server, ffmpegPath := gatherVars()
	conn, err := net.Dial("udp", server)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer conn.Close()

	cmd := exec.Command(ffmpegPath,
		"-f", "dshow",
		"-i", "audio=CABLE Output (VB-Audio Virtual Cable)",
		"-ac", "2",
		"-ar", "48000",
		"-c:a", "libopus",
		"-b:a", "64k",
		"-frame_size", "960",
		"-f", "rtp",
		"rtp:"+server,
	)

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start FFmpeg: %v", err)
	}

	log.Println("Streaming audio to bot...")
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("FFmpeg error: %v", err)
	}

}

func gatherVars() (string, string) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %s", err)
	}
	port := bot.GetEnvOrDefault("AUDIO_SERVER_PORT", "5005")
	server := bot.GetEnvOrDefault("AUDIO_SERVER_IP", "127.0.0.1")
	ffmpegPath := bot.GetEnvOrDefault("FFMPEG_PATH", "ffmpeg")

	addr := fmt.Sprintf("%s:%s", server, port)
	return addr, ffmpegPath
}
