package bot

import (
	"fmt"
	"log"
	"strings"
	"time"
	"venova/db"

	"github.com/bwmarrin/discordgo"
)

var tcGeneralId string = "209403061205073931"
var morthisId string = "186317976033558528"
var bettyId string = "641009995634180096"
var venovaId string = "1163950982259036302"
var blueId string = "202213189482446851"
var channelId string = "209404729225248769"
var griefers []string = []string{}
var nick string

func OnReady(discord *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s\n", event.User.String())
}

func HandleMessageEvents(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == discord.State.User.ID {
		return
	}

	AddGriefer(discord, msg)
	HandleCommands(discord, msg)
}

func AddGriefer(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	parts := strings.Split(msg.Content, " ")

	if parts[0] == "!grief" {
		if len(msg.Mentions) == 0 {
			if len(griefers) == 0 {
				discord.ChannelMessageSend(msg.ChannelID, "Nobody is getting griefed!")

				return
			} else {
				myGriefees := []string{}

				for _, grief := range griefers {
					myGriefees = append(myGriefees, fmt.Sprintf("<@%s>", grief))
				}
				discord.ChannelMessageSend(msg.ChannelID, strings.Join(myGriefees, " "))

				return
			}
		}

		for _, mention := range msg.Mentions {
			griefers = append(griefers, mention.ID)
		}

		discord.ChannelMessageSend(msg.ChannelID, "This brotha is getting griefed")

		return
	}
}

func HandleVoiceStateUpdate(discord *discordgo.Session, msg *discordgo.VoiceStateUpdate) {
	if msg.ChannelID != channelId {
		return
	}

	if msg.VoiceState.UserID == morthisId {
		discord.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Hello cutie <@%s>", morthisId))
	}

	for _, griefee := range griefers {
		if msg.VoiceState.UserID == griefee {
			discord.GuildMemberMove(msg.GuildID, griefee, &channelId)
		}
	}
}

func birthdateCheck(discord *discordgo.Session) {
	currDate := time.Now()

	birthDateDiscId, err := db.GetBirthdays(currDate)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for _, value := range birthDateDiscId {
		discord.ChannelMessageSend(tcGeneralId, fmt.Sprintf("Happy Birthday <@%d>", value))
		dmChannel, err := discord.UserChannelCreate(bettyId)
		if err != nil {
			log.Println("Error: ", err)
			return
		}
		discord.ChannelMessageSend(dmChannel.ID, fmt.Sprintf("It's <@%d>'s birthday!", value))
	}
}
func BirthdateCheckRoutine(discord *discordgo.Session) {
	birthdateCheck(discord)
	timer := time.NewTicker(24 * time.Hour)
	for range timer.C {
		birthdateCheck(discord)
	}

}
