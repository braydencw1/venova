package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/braydencw1/venova"
	"github.com/joho/godotenv"
)

func main() {
	showVersion := flag.Bool("version", false, "Display version information")
	showVersionShort := flag.Bool("v", false, "Display version information (short flag)")
	flag.Parse()
	if *showVersion || *showVersionShort {
		ver := venova.GetVersion("venova-audio")
		// fmt.Println(ver)
		fmt.Printf("Name:\t\t%s\nVersion:\t%s\nGit revision:\t%s\nGit ref:\t%s\nGO version:\t%s\nBuilt:\t\t%s\nOS/Arch:\t%s/%s\n",
			ver.Name, ver.Version, ver.Revision, ver.Reference, ver.GoVers, ver.BuiltAt, ver.OS, ver.Arch)
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
		"rtp:"+server, // Dynamically append server IP and port
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

func gatherVars() (string, string) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %s", err)
	}
	port := os.Getenv("AUDIO_SERVER_PORT")
	if port == "" {
		log.Printf("AUDIO_SERVER_PORT not defined, defaulting to 5005")
		port = "5005"
	}

	server := os.Getenv("AUDIO_SERVER_IP")
	if server == "" {
		log.Printf("Defaulting AUDIO_SERVER_IP to 127.0.0.1.")
		server = "127.0.0.1"
	}

	ffmpegPath := os.Getenv("FFMPEG_PATH")
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg" // Default to system PATH
	}

	addr := fmt.Sprintf("%s:%s", server, port)
	return addr, ffmpegPath
}
