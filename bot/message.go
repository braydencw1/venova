package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func sendChannelMsg(discord *discordgo.Session, channelId, msg string) {
	_, err := discord.ChannelMessageSend(channelId, msg)
	if err != nil {
		log.Printf("err msgsend sendChannelMsg %s", err)
	}
}

func dmUser(discord *discordgo.Session, userId, msg string) {
	dmChannel, err := discord.UserChannelCreate(userId)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	_, err = discord.ChannelMessageSend(dmChannel.ID, msg)
	if err != nil {
		log.Printf("error msg send dmUser %s", err)
	}
}
