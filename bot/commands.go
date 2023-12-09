package bot

import (
	"fmt"
	"log"
	"strings"
	"venova/db"

	"github.com/bwmarrin/discordgo"
)

func HandleCommands(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	nick = msg.Author.Username
	if msg.Author.ID != discord.State.User.ID {
		log.Printf(msg.Author.Username + ": " + msg.Content)
	}

	if msg.Content == "!hello" {
		discord.ChannelMessageSend(msg.ChannelID, "Hello, "+msg.Author.Username+"!")
	}
	if msg.Content == fmt.Sprintf("<@%v>", venovaId) {
		discord.ChannelMessageSend(msg.ChannelID, strings.ReplaceAll(db.DndMsgResponse(), "{nick}", nick))
	}
	if msg.Content == "https://imgur.com/a/XQ3pPTQ" {
		discord.ChannelMessageSend(msg.ChannelID, "Assemble!!!!!")
	}
	if msg.Content == "!restart" && (msg.Author.ID == morthisId || msg.Author.ID == blueId || msg.Author.ID == adonId) {
		discord.ChannelMessageSend(msg.ChannelID, "Restarting the Minecraft Server. Please allow 3 minutes or so for it to come back online.")
		restartMinecraft()
	}
	if strings.HasPrefix(msg.Content, "!mc") && (msg.Author.ID == morthisId || msg.Author.ID == blueId || msg.Author.ID == adonId) {
		log.Printf("Initial Minecraft Command, %s", msg.Content)
		cmdWord := "!mc"
		stripped := stripCommand(msg.Content, cmdWord)
		log.Printf("Stripped Minecraft Command, %s", stripped)
		res, err := minecraftCommand(stripped)
		if err != nil {
			log.Printf("Err: %s", err)
		}
		discord.ChannelMessageSend(msg.ChannelID, res)

	}
	if strings.HasPrefix(msg.Content, "!whitelist") {
		cmdWord := "!whitelist"
		stripped := stripCommand(msg.Content, cmdWord)
		res, err := minecraftCommand(fmt.Sprintf("whitelist add %s", stripped))
		if err != nil {
			log.Printf("Err: %s", err)
		}
		discord.ChannelMessageSend(msg.ChannelID, res)

	}
}
func stripCommand(input string, cmdWord string) string {
	stripped := strings.Replace(input, cmdWord, "", 1)
	stripped = strings.TrimSpace(stripped)
	return stripped
}
