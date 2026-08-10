package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/braydencw1/venova/db"

	"github.com/bwmarrin/discordgo"
)

// InitBotID reads BOT_ID from the environment. It must be called once,
// after .env has been loaded, before any handler that references venovaId runs.
func InitBotID() {
	venovaId = os.Getenv("BOT_ID")
	if venovaId == "" {
		log.Fatalf("Must provide BOT_ID")
	}
}

var venovaId string

var bangersRoleId string = "1079585245575270480"

func OnReady(discord *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s\n", event.User.String())
}

func helpCmd(ctx CommandCtx) error {
	commands := InitCommands()
	if len(ctx.Args) == 0 {
		cmdNames := commands.ListCommandsFor(ctx)
		return ctx.Reply(fmt.Sprintf("Available commands: %s", strings.Join(cmdNames, ", ")))
	}

	// Commands the caller can't use are hidden, not explained.
	cmd, exists := commands.commands[strings.ToLower(ctx.Args[0])]
	if exists && cmd.perm.Allows(ctx) {
		return ctx.Reply(cmd.help)
	}
	return ctx.Reply("Unknown command. Use !help to see all available commands.")
}

func HandleMessageEvents(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == discord.State.User.ID || msg.Author.Bot {
		return
	}
	log.Printf("%s: %s", msg.Author.Username, msg.Content)

	if msg.Content == fmt.Sprintf("<@%v>", venovaId) {
		_, err := discord.ChannelMessageSend(msg.ChannelID, strings.ReplaceAll(db.DndMsgResponse(), "{nick}", msg.Author.Username))
		if err != nil {
			log.Printf("error sending message inside HandleMessageEvents: %s", err)
		}
	} else if msg.Content == fmt.Sprintf("<@&%s>", bangersRoleId) {
		_, err := discord.ChannelMessageSend(msg.ChannelID, "https://imgur.com/K7lTDGU")
		if err != nil {
			log.Printf("error sending message inside HandleMessageEvents: %s", err)
		}
	}
}

func GetUsernameFromID(session *discordgo.Session, userID string) (string, error) {
	user, err := session.User(userID)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func PlayDateCheckRoutine(discord *discordgo.Session) {
	playDateCheck(discord)
	timer := time.NewTicker(24 * time.Hour)
	for range timer.C {
		playDateCheck(discord)
	}
}

func playDateCheck(discord *discordgo.Session) {
	nextDay := time.Now().Add(24 * time.Hour)
	res, tcId, roleId, err := db.GetPlayDates(nextDay)
	if err != nil {
		log.Printf("Failed to get play dates: %v", err)
		return
	}

	msg := fmt.Sprintf("Dnd is scheduled for tomorrow <@&%v>", roleId)
	if res {
		_, err := discord.ChannelMessageSend(fmt.Sprintf("%v", tcId), msg)
		if err != nil {
			log.Printf("err send msg palyDateCheck %s", err)
		}
	}
}

