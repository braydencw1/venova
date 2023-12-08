package bot

import (
	"log"
	"venova/sshd"
)

func restartMinecraft() {
	client, err := sshd.ConnectToDev()
	if err != nil {
		log.Printf("SSH error: %v", err)
	}
	com, err := sshd.RunCommand(client, "docker-compose -f /app/docker-compose.yml restart")
	if err != nil {
		log.Println("Error:", err.Error())
	}
	log.Printf("Compose command response: %s", com)
}
