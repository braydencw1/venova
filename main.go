package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"venova/db"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// var tcTrainingId string = "1026643518942355476"
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
	dbUsername := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbDB := os.Getenv("DB_DB")
	dbPassword := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUsername, dbPassword, dbHost, dbPort, dbDB)
	err = db.OpenDatabase(dsn)
	if err != nil {
		log.Panicf("Database connection is rough, to say the least: %v", err)
	}
	//	birthdateCheck(discord)
	go birthdateCheckRoutine(discord)
	select {} // Block the main goroutine indefinitely
}

func onReady(discord *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s\n", event.User.String())
}

func handleMessageEvents(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == discord.State.User.ID {
		return
	}

	addGriefer(discord, msg)
	handleCommands(discord, msg)
}

func addGriefer(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	parts := strings.Split(msg.Content, " ")

	if parts[0] == "!grief" {
		if len(msg.Mentions) == 0 {
			if len(griefers) == 0 {
				discord.ChannelMessageSend(msg.ChannelID, "Nobody is getting griefed!")

				return
			} else {
				myGriefees := []string{}

				for _, grief := range griefers {
					myGriefees = append(myGriefees, fmt.Sprintf("<@%s>", grief))
				}
				discord.ChannelMessageSend(msg.ChannelID, strings.Join(myGriefees, " "))

				return
			}
		}

		for _, mention := range msg.Mentions {
			griefers = append(griefers, mention.ID)
		}

		discord.ChannelMessageSend(msg.ChannelID, "This brotha is getting griefed")

		return
	}
}

func handleCommands(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID != discord.State.User.ID {
		log.Printf(msg.Author.Username + ": " + msg.Content)
	}

	if msg.Content == "!hello" {
		discord.ChannelMessageSend(msg.ChannelID, "Hello, "+msg.Author.Username+"!")
	}

	if msg.Content == fmt.Sprintf("<@%v>", vetroId) {
		discord.ChannelMessageSend(msg.ChannelID, "Hey, "+msg.Author.Username+"!")

	}
	if msg.Content == "https://imgur.com/a/XQ3pPTQ" {
		discord.ChannelMessageSend(msg.ChannelID, "Assemble!!!!!")
	}
}

func handleVoiceStateUpdate(discord *discordgo.Session, msg *discordgo.VoiceStateUpdate) {
	if msg.ChannelID != channelId {
		return
	}

	if msg.VoiceState.UserID == morthisId {
		discord.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Hello gaylord <@%s>", morthisId))
	}

	for _, griefee := range griefers {
		if msg.VoiceState.UserID == griefee {
			discord.GuildMemberMove(msg.GuildID, griefee, &channelId)
		}
	}
}

func birthdateCheck(discord *discordgo.Session) {
	currDate := time.Now()

	birthDateDiscId, err := db.GetBirthdays(currDate)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for _, value := range birthDateDiscId {
		discord.ChannelMessageSend(tcGeneralId, fmt.Sprintf("Happy Birthday <@%d>", value))
	}
}
func birthdateCheckRoutine(discord *discordgo.Session) {
	birthdateCheck(discord)
	timer := time.NewTicker(24 * time.Hour)
	for range timer.C {
		birthdateCheck(discord)
	}

}
