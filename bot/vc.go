package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func HandleAFK(discord *discordgo.Session, vc *discordgo.VoiceStateUpdate) {
	if vc.VoiceState.ChannelID != "" {
		log.Printf("Test, %v", vc.UserID)
	}
}
