package bot

import (
	"fmt"
	"log"
	"time"
	"venova/db"

	"github.com/bwmarrin/discordgo"
)

func birthdateCheck(discord *discordgo.Session) {
	nextDay := time.Now()

	birthDateDiscId, err := db.GetBirthdays(nextDay)
	if err != nil {
		fmt.Println("Error: ", err)
		log.Printf("HEREEE")
		return
	}
	for _, value := range birthDateDiscId {
		sendChannelMsg(discord, tcGeneralId, fmt.Sprintf("Happy Birthday <@%d>", value))
		dmUser(discord, bettyId, fmt.Sprintf("It's <@%d>'s birthday!", value))

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
