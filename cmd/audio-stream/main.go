package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	server := gatherVars()
	conn, err := net.Dial("udp", server)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer conn.Close()

	cmd := exec.Command("ffmpeg",
		"-f", "dshow", // Capture from DirectShow (Windows)
		"-i", "audio=CABLE Output (VB-Audio Virtual Cable)", // Adjust device name
		"-ac", "2", "-ar", "48000", "-b:a", "96k", // Stereo, 48kHz, 96kbps Opus
		"-c:a", "libopus", // Encode in Opus format
		"-f", "opus", "udp://"+server, // Stream via UDP
	)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start FFmpeg: %v", err)
	}

	log.Println("Streaming audio to bot...")
	err = cmd.Wait() // Blocks and waits for ffmpeg to finish
	if err != nil {
		log.Fatalf("FFmpeg error: %v", err)
	}

}

func gatherVars() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	port := os.Getenv("AUDIO_SERVER_PORTBOT_SERVER")
	if port == "" {
		log.Printf("AUDIO_SERVER_PORT not defined, defaulting to 5005")
		port = "5005"
	}
	server := os.Getenv("BOT_SERVER_IP")
	if server == "" {
		log.Fatalf("BOT_SERVER_IP not defined. Exiting.")
	}
	addr := fmt.Sprintf("%s:%s", server, port)
	return addr
}
