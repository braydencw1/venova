package bot

import (
	"fmt"
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
			return ctx.Reply("Your dnd role is not found in the db.")
		}
		if err := db.InsertPlayDate(t, currRoleId); err != nil {
			return fmt.Errorf("error inserting into table %w", err)
		}
		if err = ctx.Reply("The Date has been updated."); err != nil {
			return err
		}
	}
	return nil
}

func whenIsDndCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	now := time.Now()
	currRoleId := getMemberDNDRole(msg.Member)
	if currRoleId == "" {
		return ctx.Reply("Could not find DND role")
	}
	dateOfPlay, _, err := db.GetLatestPlayDate(currRoleId)
	if err != nil {
		return ctx.Reply("Could not find play date information. Perhaps wrong server.")
	}
	fmtDate := fmt.Sprint(dateOfPlay.Format("01-02-2006"))
	if dateOfPlay.Before(now) {
		if err := ctx.Reply(fmt.Sprintf("There is no date currently set. Your last session was: %s", fmtDate)); err != nil {
			return err
		}
	}
	if err = ctx.Reply(fmtDate); err != nil {
		return err
	}
	return nil
}
