package bot

import (
	"fmt"
	"log"
	"time"
	"venova/db"
)

func playCmd(ctx CommandCtx) error {
	msg := ctx.Message
	args := ctx.Args
	// Set dnd play date
	if msg.Author.ID == morthisId {
		layout := "01-02-2006"
		t, err := time.Parse(layout, args[1])
		if err != nil {
			return fmt.Errorf("error parsing date: %w", err)
		}
		currRoleId := getMemberDNDRole(msg.Member)
		if currRoleId == "" {
			log.Printf("Role not found.")
			ctx.Reply("Your dnd role is not found in the db.")
		} else {
			err := db.InsertPlayDate(t, currRoleId)
			if err != nil {
				return fmt.Errorf("error inserting into table %w", err)
			}
			ctx.Reply("The Date has been updated.")
			return nil
		}
	}
	return nil
}
func whenIsDndCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	now := time.Now()
	currRoleId := getMemberDNDRole(msg.Member)
	if currRoleId == "" {
		log.Printf("Could not find Dnd Role")
	}
	dateOfPlay, _, err := db.GetLatestPlayDate(currRoleId)
	if err != nil {
		ctx.Reply("Could not find play date information. Perhaps wrong server.")
		return fmt.Errorf("error parsing latest playdate %w", err)
	}
	fmtDate := fmt.Sprint(dateOfPlay.Format("01-02-2006"))
	if dateOfPlay.Before(now) {
		ctx.Reply(fmt.Sprintf("There is no date currently set. Your last session was: %s", fmtDate))
	} else {
		ctx.Reply(fmtDate)
	}
	return nil
}
