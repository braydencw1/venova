package main

import (
	"fmt"
	"log"
	"os"
	"venova/bot"
	"venova/db"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

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
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentGuildMembers | discordgo.IntentGuildPresences

	// Open a connection to Discord
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord connection:", err)
	}
	defer discord.Close()

	discord.AddHandler(bot.OnReady)
	discord.AddHandler(bot.HandleMessageEvents)
	discord.AddHandler(bot.HandleVoiceStateUpdate)

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
	go bot.BirthdateCheckRoutine(discord)
	go bot.PlayDateCheckRoutine(discord)
	select {} // Block the main goroutine indefinitely
}
