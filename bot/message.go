package bot

import (
	"log"

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
