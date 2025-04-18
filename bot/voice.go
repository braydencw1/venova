package bot

import (
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var activeReceiver *AudioReceiver

func stopActiveReceiver() {
	if activeReceiver != nil {
		activeReceiver.Stop()
		activeReceiver = nil
	}
}

func dcCmd(ctx CommandCtx) error {
	if !ctx.IDChecker.IsAdmin(ctx.Message.Author.ID) {
		return ctx.Reply("Cannot use this cmd.")
	}

	var uId string
	if len(ctx.Args) < 1 {
		uId = venovaId
	} else {
		uId = strings.Trim(ctx.Args[0], "<@>")
	}

	if uId == venovaId {
		stopActiveReceiver()
	}
	err := ctx.Session.GuildMemberMove(ctx.Message.GuildID, uId, nil)
	if err != nil {
		return ctx.Reply("Error disconnecting from voice channel.")
	}

	if uId == venovaId {
		return ctx.Reply("Disconnected Venova from the server.")
	}
	return ctx.Reply("Disconnected user from voice channel.")
}

func monitorVoiceActivity(session *discordgo.Session, guildID, channelID string) {
	ticker := time.NewTicker(1 * time.Minute)

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
		if userCount != 1 {
			return
		}
		log.Println("Bot is alone in VC, disconnecting...")
		err = LeaveVoiceChannel(session, guildID)
		if err != nil {
			log.Printf("Error disconnecting bot: %v", err)
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
