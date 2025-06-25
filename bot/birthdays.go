package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/braydencw1/venova/db"

	"github.com/bwmarrin/discordgo"
)

func birthdateCheck(discord *discordgo.Session) {
	nextDay := time.Now()

	bdayMessages, err := db.GetBirthdays(nextDay)
	if err != nil {
		log.Printf("error fetching birthdates :%s", err)
		return
	}

	for _, bdayMsg := range bdayMessages {
		response := bdayMsg.BdayResponse
		if response == "" {
			response = fmt.Sprintf("Happy Birthday <@%d>", bdayMsg.DiscordId)
		}
		sendChannelMsg(discord, bdayMsg.TextChannelID, response)

		// Reminder users who want individual reminders
		res, err := GetIdentityChecker().WantsBirthdayReminder()

		if err != nil {
			log.Printf("Error extracting Birthday Reminder Users %s", err)
			continue
		}

		for _, id := range res {
			dmUser(discord, id, fmt.Sprintf("It's <@%d>'s birthday!", bdayMsg.DiscordId))
		}
	}
}

func BirthdateCheckRoutine(discord *discordgo.Session) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}
		durTilNextCheck := next.Sub(now)
		timer := time.NewTimer(durTilNextCheck)
		<-timer.C
		birthdateCheck(discord)
	}
}
