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
	if msg.Content == "!mc_restart" && (msg.Author.ID == morthisId || msg.Author.ID == blueId) {
		discord.ChannelMessageSend(msg.ChannelID, "Restarting MWINECWAFT")
		restartMinecraft()
	}
}
