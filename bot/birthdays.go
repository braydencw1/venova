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

	birthDateDiscId, err := db.GetBirthdays(nextDay)
	if err != nil {
		log.Printf("error fetching birthdates :%s", err)
		return
	}
	for discID, bdres := range birthDateDiscId {
		log.Printf("Here is %d, %s", discID, bdres)
		if bdres == "" {
			sendChannelMsg(discord, tcGeneralId, fmt.Sprintf("Happy Birthday <@%d>", discID))
		} else {
			sendChannelMsg(discord, tcGeneralId, fmt.Sprintf("%s <@%d>", bdres, discID))
		}

		dmUser(discord, bettyId, fmt.Sprintf("It's <@%d>'s birthday!", discID))

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
