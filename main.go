package main

import (
	"log"
	"os"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)
var morthisId string = ""
var vetroId string = ""


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
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	// Open a connection to Discord
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord connection:", err)
	}
	defer discord.Close()
	discord.AddHandler(onReady)
	// Register the messageCreate functfunc messageCreate(s *discordgo.Session, m *discordgo.MessageCreate)ion as a callback for the MessageCreate event
	// discord.AddMessageCreateHandler(messageCreate)
	discord.AddHandler(messageCreate)
	// Keep the bot running
	log.Println("Bot is now running. Press Ctrl+C to exit.")
	select {} // Block the main goroutine indefinitely
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s\n", event.User.String())
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.ID != s.State.User.ID {
		log.Printf( m.Author.Username +": " + m.Content)
	}

	// Respond to messages
	if m.Content == "!hello" {
		// Reply with a message
		s.ChannelMessageSend(m.ChannelID, "Hello, "+m.Author.Username+"!")
	}

	if m.Content == fmt.Sprintf("<@%v>", vetroId) {
		s.ChannelMessageSend(m.ChannelID, "Hey, "+m.Author.Username+"!")
		
	}
	if m.Content == "https://imgur.com/a/XQ3pPTQ" {
		s.ChannelMessageSend(m.ChannelID, "Assemble!!!!!")
	}
}
