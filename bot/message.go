package bot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func sendChannelMsg(discord *discordgo.Session, channelId, msg string) {
	discord.ChannelMessageSend(channelId, msg)
}

func dmUser(discord *discordgo.Session, userId, msg string) {
	dmChannel, err := discord.UserChannelCreate(userId)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	discord.ChannelMessageSend(dmChannel.ID, msg)
}

func handleCommands(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID != discord.State.User.ID {
		log.Printf(msg.Author.Username + ": " + msg.Content)
	}

	parts := strings.SplitN(msg.Content, " ", 2)

	command := strings.TrimPrefix(parts[0], "!")

	if fn, ok := botCommandsWithArgs[command]; ok {
		if len(parts) < 2 {
			discord.ChannelMessageSend(msg.ChannelID, "Need more arguments for command.")
			return
		}
		fn(discord, msg, parts)
	} else if fn, ok := botCommandsWithoutArgs[command]; ok {
		fn(discord, msg)
	} else {
		log.Printf("Invalid Command.")
		discord.ChannelMessageSend(msg.ChannelID, "Invalid command.")
	}
}
