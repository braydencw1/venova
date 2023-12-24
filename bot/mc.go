package bot

import (
	"fmt"
	"log"
	"os"
	"venova/sshcmd"

	"github.com/gorcon/rcon"
)

func restartMinecraft() {
	client, err := sshcmd.ConnectToDev()
	if err != nil {
		log.Printf("SSH error: %v", err)
		return
	}
	defer client.Close()

	com, err := sshcmd.RunCommand(client, "docker-compose -f /app/docker-compose.yml restart")
	if err != nil {
		log.Println("Error:", err.Error())
	}
	log.Printf("Compose command response: %s", com)
}

func minecraftCommand(command string) (string, error) {

	rconHost := os.Getenv("MC_HOST")
	rconPort := os.Getenv("RCON_PORT")
	rconPass := os.Getenv("RCON_PASS")

	con, err := rcon.Dial(fmt.Sprintf("%s:%s", rconHost, rconPort), rconPass)
	if err != nil {
		log.Printf("Error, %s:", err)
		return "", err
	}
	defer con.Close()
	response, err := con.Execute(command)
	if err != nil {
		log.Printf("Error, %s", err)
		return "", err
	}
	log.Printf("Response %s", response)
	return response, nil
}
