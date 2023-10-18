package main

import (
	"fmt"
	"log"
	"time"
	"os"
	"strings"
	"venova/db"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)
var tcTrainingId string = "1026643518942355476"
var tcGeneralId string = "209403061205073931"
var morthisId string = "186317976033558528"
var vetroId string = "1131832403581747381"
var channelId string = "209404729225248769"
var griefers []string = []string{}

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	token := os.Getenv("TOKEN")

	// Initialize Discord session
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a connection to Discord
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord connection:", err)
	}
	defer discord.Close()
	discord.AddHandler(onReady)
	// Register the messageCreate functfunc messageCreate(s *discordgo.Session, m *discordgo.MessageCreate)ion as a callback for the MessageCreate event
	// discord.AddMessageCreateHandler(messageCreate)
	discord.AddHandler(handleMessageEvents)
	discord.AddHandler(handleVoiceStateUpdate)
	// Keep the bot running
	log.Println("Bot is now running. Press Ctrl+C to exit.")
	birthdateCheck(discord)
	timer := time.NewTicker(24 * time.Hour)
	go func() {
		for range timer.C {
			now = time.Now()
			birthdateCheck(discord)
		}

	}()
	select {} // Block the main goroutine indefinitely
}

func onReady(sess *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s\n", event.User.String())
}

func handleMessageEvents(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	if mess.Author.ID == sess.State.User.ID {
		return
	}

	addGriefer(sess, mess)
	handleCommands(sess, mess)
}

func addGriefer(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	parts := strings.Split(mess.Content, " ")

	if parts[0] == "!grief" {
		if len(mess.Mentions) == 0 {
			if len(griefers) == 0 {
				sess.ChannelMessageSend(mess.ChannelID, "Nobody is getting griefed!")

				return
			} else {
				myGriefees := []string{}

				for _, grief := range griefers {
					myGriefees = append(myGriefees, fmt.Sprintf("<@%s>", grief))
				}
				sess.ChannelMessageSend(mess.ChannelID, strings.Join(myGriefees, " "))

				return
			}
		}

		for _, mention := range mess.Mentions {
			griefers = append(griefers, mention.ID)
		}

		sess.ChannelMessageSend(mess.ChannelID, "This brotha is getting griefed")

		return
	}
}

func handleCommands(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	if mess.Author.ID != sess.State.User.ID {
		log.Printf(mess.Author.Username + ": " + mess.Content)
	}

	// Respond to messages
	if mess.Content == "!hello" {
		// Reply with a message
		sess.ChannelMessageSend(mess.ChannelID, "Hello, "+mess.Author.Username+"!")
	}

	if mess.Content == fmt.Sprintf("<@%v>", vetroId) {
		sess.ChannelMessageSend(mess.ChannelID, "Hey, "+mess.Author.Username+"!")

	}
	if mess.Content == "https://imgur.com/a/XQ3pPTQ" {
		sess.ChannelMessageSend(mess.ChannelID, "Assemble!!!!!")
	}
}

func handleVoiceStateUpdate(sess *discordgo.Session, mess *discordgo.VoiceStateUpdate) {
	if mess.ChannelID != channelId {
		return
	}

	if mess.VoiceState.UserID == morthisId {
		sess.ChannelMessageSend(mess.ChannelID, fmt.Sprintf("Hello gaylord <@%s>", morthisId))
	}

	for _, griefee := range griefers {
		if mess.VoiceState.UserID == griefee {
			sess.GuildMemberMove(mess.GuildID, griefee, &channelId)
		}
	}
}

func birthdateCheck(sess *discordgo.Session) { 
	var birthDateDiscId int
	now := time.Now()
	currDate := fmt.Sprintf(now.Format("2006-01-02"))
        targetTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location()) 
        if now.Hour() == targetTime.Hour() { 
		birthDateDiscId = birthdateQuery(currDate)
	}
	birthDateDiscId = birthdateQuery(currDate)
	if birthDateDiscId != -1 {
		birthdayMessage := fmt.Sprintf("Happy Birthday: <@%d>", birthDateDiscId)

		_, err := sess.ChannelMessageSend(tcTrainingId, birthdayMessage)
		if err != nil {
			fmt.Println("Error:", err)
		
		}
	
	
	}
}


func birthdateQuery(dateToCheck string) int {
	var discIdReturn int
	myQuery := fmt.Sprintf("SELECT * FROM users WHERE TO_CHAR(dob, 'YYYY-MM-DD') LIKE '%%%s%%';", dateToCheck)
        _, rows, err := db.Dquery(myQuery) 
        if err != nil { 
                fmt.Println("Error:", err) 
		return -1
        }
	defer rows.Close()
	var firstname, lastname, initbday, ignored string
	var discid int
	for rows.Next() {
		err := rows.Scan(&ignored, &discid, &firstname, &lastname, &initbday)
		if err != nil {
		fmt.Println("Error:", err)
			return -1
		}
		discIdReturn = discid
 
	}
	return discIdReturn
}
