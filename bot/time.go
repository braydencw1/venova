package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/braydencw1/venova/db"
	"github.com/bwmarrin/discordgo"
)

func setTimerCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	if !ctx.IDChecker.IsAdmin(msg.Author.ID) {
		return nil
	}

	duration, err := time.ParseDuration(ctx.Args[0])
	if err != nil {
		return ctx.Reply(fmt.Sprintf("Invalid duration %q. Use e.g. 30m or 1h30m.", ctx.Args[0]))
	}

	targetID := msg.Author.ID
	if len(ctx.Args) > 1 {
		// Accept both <@id> and the legacy nickname form <@!id>.
		targetID = strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(ctx.Args[1], "<@"), "!"), ">")
	}

	timer := db.Timer{
		TargetID:    targetID,
		CreatedByID: msg.Author.ID,
		FiresAt:     time.Now().Add(duration),
	}
	if err := db.InsertTimer(&timer); err != nil {
		return fmt.Errorf("saving timer: %w", err)
	}
	log.Printf("Timer %d set by %s for %s, firing at %s", timer.ID, msg.Author.Username, targetID, timer.FiresAt)

	go runTimer(ctx.Session, timer)

	return ctx.Reply(fmt.Sprintf("Timer set for %s — I'll DM <@%s> when it's up.", duration, targetID))
}

// runTimer waits until the timer's fire time, DMs the target, and removes the
// timer from the database. The row is only deleted after a successful send so
// a crash or failed DM gets retried on the next startup.
func runTimer(s *discordgo.Session, t db.Timer) {
	time.Sleep(time.Until(t.FiresAt))

	dmChannel, err := s.UserChannelCreate(t.TargetID)
	if err != nil {
		log.Printf("timer %d: creating dm channel for %s: %s", t.ID, t.TargetID, err)
		return
	}
	if _, err := s.ChannelMessageSend(dmChannel.ID, "Your timer is up!"); err != nil {
		log.Printf("timer %d: sending dm to %s: %s", t.ID, t.TargetID, err)
		return
	}
	if err := db.DeleteTimer(t.ID); err != nil {
		log.Printf("timer %d: deleting fired timer: %s", t.ID, err)
	}
}

// ResumeTimers reschedules every timer still in the database. Timers whose
// fire time passed while the bot was down fire immediately.
func ResumeTimers(s *discordgo.Session) {
	timers, err := db.GetPendingTimers()
	if err != nil {
		log.Printf("resuming timers: %s", err)
		return
	}
	for _, t := range timers {
		go runTimer(s, t)
	}
	if len(timers) > 0 {
		log.Printf("resumed %d timer(s)", len(timers))
	}
}
