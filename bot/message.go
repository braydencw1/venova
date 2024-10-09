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

// func HandleCommands(discord *discordgo.Session, msg *discordgo.MessageCreate) {
// 	if msg.Author.ID != discord.State.User.ID {
// 		log.Printf(msg.Author.Username + ": " + msg.Content)
// 	}

// 	parts := strings.SplitN(msg.Content, " ", 2)
// 	guild, err := discord.State.Guild(msg.GuildID)
// 	if err != nil {
// 		log.Printf("Failed to fetch message guild id (dm?): %v", err)
// 		return
// 	}

// 	member := getGuildMember(guild, msg.Author.ID)
// 	if member == nil {
// 		log.Printf("Failed to find message guild member: %v", msg.Author.ID)
// 		return
// 	}

// 	if fn, ok := BotCommands[fmt.Sprintf("!%s", parts[0])]; ok {
// 		fn(discord, msg, parts)
// 	} else {
// 		log.Printf("Invalid Command.")
// 		discord.ChannelMessageSend(msg.ChannelID, "Invalid command.", nil)
// 	}
// }
