package bot

import (
	"fmt"
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

func HandleCommands(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID != discord.State.User.ID {
		log.Printf(msg.Author.Username + ": " + msg.Content)
	}

	parts := strings.SplitN(msg.Content, " ", 2)

	command := strings.TrimPrefix(parts[0], "!")
	if len(parts) < 2 {
		if fn, ok := botCommandsWithArgs[fmt.Sprintf(command)]; ok {
			fn(discord, msg, parts)
		} else {
			log.Printf("Invalid Command.")
			discord.ChannelMessageSend(msg.ChannelID, "Invalid command.")
		}
	}
}
