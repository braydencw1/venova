package bot

import (
	"fmt"
	"log"
	"os"
	"time"
	"venova/sshcmd"

	"github.com/gorcon/rcon"
)

func manageMinecraftCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	sess := ctx.Session
	args := ctx.Args
	actions := map[string]struct {
		Message string
		Command string
	}{
		"up":      {"The Minecraft Server has been brought up", "up -d"},
		"down":    {"The Minecraft server has been brought down", "down"},
		"restart": {"The Minecraft server has been restarted.", "restart"},
	}
	if action, exists := actions[args[0]]; exists {
		if memberHasRole(msg.Member, mcRoleId) || msg.Author.ID == morthisId {
			mcMsg, _ := sess.ChannelMessageSend(msg.ChannelID, "Attempting to modify the minecraft server...")

			go func() {
				err := execDockerCompose(action.Command)
				if err != nil {
					log.Printf("%s", err)
					ctx.Reply(fmt.Sprintf("%s", err))
					return
				}
				time.Sleep(time.Second * 5)
				sess.ChannelMessageEdit(msg.ChannelID, mcMsg.ID, action.Message)
			}()
			return nil
		}
	} else {
		ctx.Reply("Action is not available")
		return nil
	}
	ctx.Reply("You're not a minecraft administrator.")
	return nil
}

func mcCmd(ctx CommandCtx) error {
	m := ctx.Message
	if m.Author.ID == blueId || m.Author.ID == morthisId {
		res, err := minecraftCommand(ctx.Args[0])
		if err != nil {
			log.Printf("Err: %s", err)
			return ctx.Reply("Could not send command, Minecraft might be offline.")
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

func minecraftCommand(command string) (string, error) {
	rconHost := os.Getenv("MC_HOST")
	rconPort := os.Getenv("RCON_PORT")
	rconPass := os.Getenv("RCON_PASS")

	if rconHost == "" || rconPort == "" || rconPass == "" {
		err := fmt.Errorf("missing rcon env vars")
		log.Printf("Error: %v", err)
		return "", err
	}

	con, err := rcon.Dial(fmt.Sprintf("%s:%s", rconHost, rconPort), rconPass)
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
	log.Printf("Response %s", resp)
	return resp, nil
}

func whitelistCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	args := ctx.Args
	if memberHasRole(msg.Member, mcRoleId) || msg.Author.ID == morthisId {
		log.Printf("Whitelisting, %s ", args[0])
		res, err := minecraftCommand(fmt.Sprintf("whitelist add %s", args[0]))
		if err != nil {
			ctx.Reply("Could not send command, Minecraft might be offline.")
			return fmt.Errorf("minecraft might be offline. err: %w", err)
		}
		ctx.Reply(res)
	} else {
		ctx.Reply("You're not a miencraft admin.")
	}
	return nil
}
