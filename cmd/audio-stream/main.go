package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	server := gatherVars()
	conn, err := net.Dial("udp", server)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("%s", conn)

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
