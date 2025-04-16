package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/braydencw1/venova/sshcmd"

	"github.com/gorcon/rcon"
)

type minecraftAction struct {
	Message string
	Command string
}

var minecraftActions = map[string]minecraftAction{
	"up":      {Message: "The Minecraft Server has been brought up", Command: "up -d"},
	"down":    {Message: "The Minecraft server has been brought down", Command: "down"},
	"restart": {Message: "The Minecraft server has been restarted.", Command: "restart"},
}

func manageMinecraftCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	sess := ctx.Session
	args := ctx.Args
	action, exists := minecraftActions[args[0]]
	if !exists {
		if err := ctx.Reply("Action is not available."); err != nil {
			return err
		}
	}
	if memberHasRole(msg.Member, mcRoleId) || msg.Author.ID == morthisId {
		mcMsg, _ := sess.ChannelMessageSend(msg.ChannelID, "Attempting to modify the minecraft server...")
		go func() {
			err := execDockerCompose(action.Command)
			if err != nil {
				log.Printf("%s", err)
				if err := ctx.Reply(fmt.Sprintf("%s", err)); err != nil {
					log.Printf("error replying manageMinecraftCmd %s", err)
				}
				return
			}
			time.Sleep(time.Second * 5)
			_, err = sess.ChannelMessageEdit(msg.ChannelID, mcMsg.ID, action.Message)
			if err != nil {
				log.Printf("error editting message inside of manageMiencraftCmd %s", err)
			}
		}()
	}
	return nil
}

func mcCmd(ctx CommandCtx) error {
	m := ctx.Message
	if m.Author.ID == blueId || m.Author.ID == morthisId {
		res, err := minecraftCommand(strings.Join(ctx.Args, " "))
		if err != nil {
			log.Printf("Err: %s", err)
			return ctx.Reply("Could not send command, Minecraft might be offline.")
		}
		if res == "" {
			return nil
		}
		return ctx.Reply(res)
	}
	return ctx.Reply("Sorry you're not a Minecraft server admin.")

}

func execDockerCompose(action string) error {
	client, err := sshcmd.ConnectToDev()
	if err != nil {
		return fmt.Errorf("ssh error: %w", err)
	}
	defer client.Close()

	com, err := sshcmd.RunCommand(client, fmt.Sprintf("docker-compose -f /app/docker-compose.yml %s", action))
	if err != nil {
		return fmt.Errorf("error running ssh command: %w", err)
	}
	log.Printf("Compose command response: %s", com)
	return nil
}
func initMC() (string, string, string, error) {
	rconHost := os.Getenv("MC_HOST")
	rconPort := os.Getenv("RCON_PORT")
	rconPass := os.Getenv("RCON_PASS")

	if rconHost == "" || rconPort == "" || rconPass == "" {
		err := fmt.Errorf("missing rcon env vars")
		log.Printf("Error: %v", err)
		return "", "", "", err
	}
	return rconHost, rconPort, rconPass, nil
}

func minecraftCommand(command string) (string, error) {
	host, port, pass, err := initMC()
	if err != nil {
		return "", err
	}
	log.Printf("Cmd here %s", command)
	con, err := rcon.Dial(fmt.Sprintf("%s:%s", host, port), pass)
	if err != nil {
		log.Printf("Error connecting to RCON: %s", err)
		return "", fmt.Errorf("unable to connect to rcon server: %s", err)
	}
	defer con.Close()
	resp, err := con.Execute(command)
	if err != nil {
		log.Printf("Error, %s", err)
		return "", err
	}
	return resp, nil
}

func whitelistCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	args := ctx.Args
	if memberHasRole(msg.Member, mcRoleId) || msg.Author.ID == morthisId {
		log.Printf("Whitelisting, %s ", args[0])
		res, err := minecraftCommand(fmt.Sprintf("whitelist add %s", args[0]))
		if err != nil {
			if err := ctx.Reply("Could not send command, Minecraft might be offline."); err != nil {
				return err
			}
		}
		if err := ctx.Reply(res); err != nil {
			return err
		}

	}
	return nil

}
