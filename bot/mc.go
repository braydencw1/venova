package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/braydencw1/venova/sshcmd"
	"golang.org/x/crypto/ssh"

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
	if len(args) == 0 {
		return ctx.Reply("Please supply an action")
	}
	action, exists := minecraftActions[args[0]]
	if !exists {
		if err := ctx.Reply("Action is not available."); err != nil {
			return err
		}
	}
	if ctx.IDChecker.IsMinecraftAdmin(msg.Author.ID) || ctx.IDChecker.IsAdmin(msg.Author.ID) {
		mcMsg, _ := sess.ChannelMessageSend(msg.ChannelID, "Attempting to modify the minecraft server...")
		go func() {
			err := execCompose(action.Command)
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
	mID := m.Author.ID
	if !ctx.IDChecker.IsMinecraftAdmin(mID) && !ctx.IDChecker.IsAdmin(mID) {
		return ctx.Reply("Sorry you're not a Minecraft server admin.")
	}
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

func getComposeCmd(client *ssh.Client) (string, error) {
	if _, err := sshcmd.RunCommand(client, "command -v docker-compose"); err == nil {
		return "docker-compose", nil
	}

	if _, err := sshcmd.RunCommand(client, "command -v podman-compose"); err == nil {
		return "podman-compose", nil
	}

	if _, err := sshcmd.RunCommand(client, "command -v docker"); err == nil {
		if _, err := sshcmd.RunCommand(client, "docker compose version"); err == nil {
			return "docker compose", nil
		}
	}

	return "", fmt.Errorf("compose binary not found")
}

func getComposePath(client *ssh.Client) (string, error) {
	defaultPath := os.Getenv("MC_COMPOSE_PATH")
	if defaultPath != "" {
		if _, err := sshcmd.RunCommand(client, fmt.Sprintf("test -f %s", defaultPath)); err == nil {
			return defaultPath, nil
		}
		log.Printf("MC_COMPOSE_PATH does not exist, continuing to defaults.")
	}

	paths := []string{
		"/app/docker-compose/docker-compose.yml",
		"/app/docker-compose.yml",
		"/app/podman-compose.yml",
		"/app/podman-compose/podman-compose.yml",
	}

	for _, p := range paths {
		if _, err := sshcmd.RunCommand(client, fmt.Sprintf("test -f %s", p)); err == nil {
			return p, nil
		}

	}

	return "", fmt.Errorf("no compose command path")
}

func execCompose(action string) error {
	client, err := sshcmd.ConnectToDev()
	if err != nil {
		return fmt.Errorf("ssh error: %w", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("ERROR: %s", err)
		}

	}()

	composeCmd, err := getComposeCmd(client)
	log.Printf("compose found: %s", composeCmd)
	if err != nil {
		return err
	}

	composePath, err := getComposePath(client)
	if err != nil {
		return err
	}

	com, err := sshcmd.RunCommand(client, fmt.Sprintf("%s -f %s %s", composeCmd, composePath, action))
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

	defer func() {
		if err := con.Close(); err != nil {
			log.Printf("ERROR: %s", err)
		}
	}()

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
	if !ctx.IDChecker.IsMinecraftAdmin(msg.Author.ID) && !ctx.IDChecker.IsAdmin(msg.Author.ID) {
		return nil
	}
	res, err := minecraftCommand(fmt.Sprintf("whitelist add %s", args[0]))
	if err != nil {
		replyErr := ctx.Reply("Could not send command, Minecraft might be offline.")
		if replyErr != nil {
			return replyErr
		}
		return nil
	}
	log.Printf("Whitelisting, %s ", args[0])
	return ctx.Reply(fmt.Sprintf("Whitelisting Res, %s", res))
}
