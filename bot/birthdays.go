package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/braydencw1/venova/db"

	"github.com/bwmarrin/discordgo"
)

// Birthdays are boundaried by the community's local calendar day, not by
// wherever the bot happens to run. Central covers both CST and CDT.
const birthdayTZ = "America/Chicago"

var birthdayLoc = mustLoadLocation(birthdayTZ)

// mustLoadLocation resolves a zone or refuses to start. cmd/venova imports
// time/tzdata, so this cannot fail in practice - if it ever does, that import
// was dropped and every date below would silently be UTC.
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Fatalf("could not load timezone %q, is time/tzdata still imported? %s", name, err)
	}
	return loc
}

func birthdateCheck(discord *discordgo.Session) (int, error) {
	today := time.Now().In(birthdayLoc)

	bdayMessages, err := db.GetBirthdays(today)
	if err != nil {
		log.Printf("error fetching birthdates :%s", err)
		return 0, err
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
	return len(bdayMessages), nil
}

// bdaySendCmd triggers the daily birthday check on demand instead of
// waiting for the 8am routine. Admins only.
func bdaySendCmd(ctx CommandCtx) error {
	if !ctx.IDChecker.IsAdmin(ctx.Message.Author.ID) {
		return nil
	}

	sent, err := birthdateCheck(ctx.Session)
	if err != nil {
		return ctx.Reply("Failed to fetch birthdays.")
	}
	if sent == 0 {
		return ctx.Reply("No birthdays today.")
	}
	return ctx.Reply(fmt.Sprintf("Sent %d birthday message(s).", sent))
}

func BirthdateCheckRoutine(discord *discordgo.Session) {
	for {
		now := time.Now().In(birthdayLoc)
		// Calendar arithmetic, not duration arithmetic: adding 24h across a DST
		// transition would land on 7am or 9am and stay drifted until restart.
		// time.Date normalizes the overflowing day and re-resolves the offset.
		next := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, birthdayLoc)
		if !now.Before(next) {
			next = time.Date(now.Year(), now.Month(), now.Day()+1, 8, 0, 0, 0, birthdayLoc)
		}
		durTilNextCheck := next.Sub(now)
		timer := time.NewTimer(durTilNextCheck)
		<-timer.C
		if _, err := birthdateCheck(discord); err != nil {
			log.Printf("birthdate check routine err: %s", err)
		}
	}
}
