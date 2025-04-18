package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/braydencw1/venova/db"

	"github.com/bwmarrin/discordgo"
)

var tcGeneralId string = "209403061205073931"

// var adminID string = "186317976033558528"
// var bettyId string = "641009995634180096"
var venovaId string = "1163950982259036302"

// var blueId string = "202213189482446851"
var bangersRoleId string = "1079585245575270480"

var mcRoleId string = "1183228947874459668"

// var channelId string = "209403061205073931"
var joinableRolesMap = map[string]string{
	"apes":     "1250598584534175784",
	"dorklock": "1282817878244200488",
	"bangers":  "1079585245575270480",
}

func OnReady(discord *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s\n", event.User.String())
}

func HandleMessageEvents(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == discord.State.User.ID {
		return
	}
	log.Printf("%s: %s", msg.Author.Username, msg.Content)

	if msg.Content == fmt.Sprintf("<@%v>", venovaId) {
		_, err := discord.ChannelMessageSend(msg.ChannelID, strings.ReplaceAll(db.DndMsgResponse(), "{nick}", msg.Author.Username))
		if err != nil {
			log.Printf("error sending message inside HandleMessageEvents: %s", err)
		}
	} else if msg.Content == fmt.Sprintf("<@&%s>", bangersRoleId) {
		_, err := discord.ChannelMessageSend(msg.ChannelID, "https://imgur.com/K7lTDGU")
		if err != nil {
			log.Printf("error sending message inside HandleMessageEvents: %s", err)
		}
	}
}

func GetUsernameFromID(session *discordgo.Session, userID string) (string, error) {
	user, err := session.User(userID)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func PlayDateCheckRoutine(discord *discordgo.Session) {
	playDateCheck(discord)
	timer := time.NewTicker(24 * time.Hour)
	for range timer.C {
		playDateCheck(discord)
	}
}

func playDateCheck(discord *discordgo.Session) {
	nextDay := time.Now().Add(24 * time.Hour)
	res, tcId, roleId, err := db.GetPlayDates(nextDay)
	if err != nil {
		log.Printf("Failed to get play dates: %v", err)
		return
	}

	msg := fmt.Sprintf("Dnd is scheduled for tomorrow <@&%v>", roleId)
	if res {
		_, err := discord.ChannelMessageSend(fmt.Sprintf("%v", tcId), msg)
		if err != nil {
			log.Printf("err send msg palyDateCheck %s", err)
		}
	}
}

func createTimer(timeLength string) (time.Time, error) {
	duration, err := time.ParseDuration(timeLength)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}, err
	}
	timer := time.Now().Add(duration)
	return timer, nil
}

func TimerCheckerRoutine(discord *discordgo.Session, timer time.Time, UserID string, errChan chan error) {
	ticker := time.NewTicker(1 * time.Minute) // Ticker to check every minute
	defer ticker.Stop()
	defer close(errChan)
	for {
		<-ticker.C
		if time.Now().After(timer) {
			dmChannel, err := discord.UserChannelCreate(UserID)
			if err != nil {
				errChan <- fmt.Errorf("error creating dm channel : %w", err)
				return
			}
			_, err = discord.ChannelMessageSend(dmChannel.ID, "Your timer is up!")
			if err != nil {
				errChan <- fmt.Errorf("error sending dm for time : %w", err)
			}
			return
		}
	}
}
