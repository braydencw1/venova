package bot

import (
	"log"
	"strings"
)

func setTimerCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	sess := ctx.Session
	args := ctx.Args
	if msg.Author.ID == morthisId || msg.Author.ID == bettyId {
		extraParts := strings.SplitN(args[0], " ", 2)
		log.Printf("Creating a timer for %s", msg.Author.Username)
		if len(extraParts) == 1 {
			extraParts = append(extraParts, msg.Author.ID)
		} else if len(extraParts) == 2 {
			result := strings.TrimPrefix(extraParts[1], "<@")
			result = strings.TrimSuffix(result, ">")
			extraParts[1] = result
		}

		timer, err := createTimer(extraParts[0])
		if err != nil {
			log.Printf("Could not create timer %s", err)
		}
		timerDestUserName, err := GetUsernameFromID(sess, extraParts[1])
		if err != nil {
			log.Printf("Could not retriever UserName from userID, %s", err)
		}
		log.Printf("Creating a timer for %s, for the length %s, destined for %s with userID %s", msg.Author.Username, extraParts[0], timerDestUserName, extraParts[1])
		if err != nil {
			log.Printf("%s", err)
		}
		errChan := make(chan error)
		go TimerCheckerRoutine(sess, timer, extraParts[1], errChan)
		err = <-errChan
		if err != nil {
			log.Printf("Error with the timer routine, %v", err)
		}
	}
	return nil
}
