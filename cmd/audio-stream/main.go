package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
	"github.com/braydencw1/venova"
	"github.com/braydencw1/venova/pkg/util"
	"github.com/joho/godotenv"
)

type AudioSender struct {
	Address     string
	FfmpegPath  string
	AudioDevice string
	AudioFormat string
}

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

	as := InitAudioSender()

	conn, err := net.Dial("udp", as.Address)
	if err != nil {
		log.Fatalf("%s", err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatalf("could not close audio connection: %s", err)
		}
	}()

	cmd := exec.Command(as.FfmpegPath,
		"-f", as.AudioFormat,
		"-i", as.AudioDevice,
		"-ac", "2",
		"-ar", "48000",
		"-c:a", "libopus",
		"-b:a", "64k",
		"-frame_size", "960",
		"-f", "rtp",
		"rtp:"+as.Address,
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

func InitAudioSender() *AudioSender {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %s", err)
	}

	port := util.GetEnvOrDefault("AUDIO_SERVER_PORT", "5005")
	server := util.GetEnvOrDefault("AUDIO_SERVER_IP", "127.0.0.1")
	ffmpegPath := util.GetEnvOrDefault("FFMPEG_PATH", "ffmpeg")
	audioFormat := util.GetEnvOrDefault("AUDIO_FORMAT", "pulse")
	audioDevice := util.GetEnvOrDefault("AUDIO_DEVICE", "default")
	addr := fmt.Sprintf("%s:%s", server, port)
	a := AudioSender{
		Address:     addr,
		FfmpegPath:  ffmpegPath,
		AudioDevice: audioDevice,
		AudioFormat: audioFormat,
	}

	log.Printf("Streaming from device %q using %s to %s", audioDevice, audioFormat, addr)

	return &a
}
