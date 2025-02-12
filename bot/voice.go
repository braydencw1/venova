package bot

import (
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func dcCmd(ctx CommandCtx) error {
	if ctx.Message.Author.ID == morthisId {
		var uId string
		log.Printf("%s", ctx.Args)
		if len(ctx.Args) < 1 {
			uId = venovaId
		} else {
			uId = strings.Trim(ctx.Args[0], "<@>")
		}
		log.Printf("%s", uId)
		err := disconnectUserFromVC(ctx.Session, ctx.Message.GuildID, uId)
		if err != nil {
			return ctx.Reply("Error disconnecting from voice channel.")
		}
		if uId == venovaId {
			return ctx.Reply("Disconnected Venova from the server.")
		} else {
			return ctx.Reply("Disconnected user from voice channel.")
		}
	}
	return ctx.Reply("stinky")
}

func disconnectUserFromVC(session *discordgo.Session, guildID, userID string) error {
	err := session.GuildMemberMove(guildID, userID, nil)
	if err != nil {
		log.Printf("Failed to disconnect user %s from VC: %v", userID, err)
		return err
	}

	log.Printf("User has been disconnected from voice.")
	return nil
}

func monitorVoiceActivity(session *discordgo.Session, guildID, channelID string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		guild, err := session.State.Guild(guildID)
		if err != nil {
			log.Printf("Error fetching guild: %v", err)
			continue
		}

		userCount := 0
		for _, vs := range guild.VoiceStates {
			if vs.ChannelID == channelID {
				userCount++
			}
		}

		if userCount == 1 {
			log.Println("Bot is alone in VC, disconnecting...")
			err := LeaveVoiceChannel(session, guildID)
			if err != nil {
				log.Printf("Error disconnecting bot: %v", err)
			}
			return
		}
	}
}
func LeaveVoiceChannel(session *discordgo.Session, guildID string) error {
	voiceConn, exists := session.VoiceConnections[guildID]
	if !exists {
		return nil // Bot is not in a VC
	}

	// Disconnect from VC
	err := voiceConn.Disconnect()
	if err != nil {
		log.Printf("Error disconnecting from voice: %v", err)
		return err
	}

	log.Println("Bot has disconnected from voice channel.")
	return nil
}
