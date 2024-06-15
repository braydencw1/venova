package bot

import (
	"fmt"
	"time"
	"venova/db"

	"github.com/bwmarrin/discordgo"
)

func birthdateCheck(discord *discordgo.Session) {
	nextDay := time.Now()

	birthDateDiscId, err := db.GetBirthdays(nextDay)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for _, value := range birthDateDiscId {
		sendChannelMsg(discord, tcGeneralId, fmt.Sprintf("Happy Birthday <@%d>", value))
		dmUser(discord, bettyId, fmt.Sprintf("It's <@%d>'s birthday!", value))

	}
}
func BirthdateCheckRoutine(discord *discordgo.Session) {
	birthdateCheck(discord)
	timer := time.NewTicker(24 * time.Hour)
	for range timer.C {
		birthdateCheck(discord)
	}

}
