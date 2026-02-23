package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/braydencw1/venova"
	"github.com/braydencw1/venova/bot"
	"github.com/braydencw1/venova/db"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var cli struct {
	Version bool `help:"Show version" short:"v"`
}
var (
	dsn   string
	token string
)

func main() {
	handleCommandLine()
	discord := startDiscordSession()
	defer func() {
		if err := discord.Close(); err != nil {
			log.Printf("failed to close Discord session: %v", err)
		}
	}()

	err := db.OpenDatabase(dsn)
	if err != nil {
		log.Panicf("Database connection is rough, to say the least: %v", err)
	}

	discord.AddHandler(bot.OnReady)
	discord.AddHandler(bot.HandleMessageEvents)

	discord.AddHandler(bot.InitCommands())

	streamers, err := db.GetStreamers()
	if err != nil {
		log.Printf("Could not initialize streamers: %w", err)
	}

	go bot.PollStreamer(discord, streamers)
	go bot.BirthdateCheckRoutine(discord)
	go bot.PlayDateCheckRoutine(discord)
	select {} // Block the main goroutine indefinitely
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dsn = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DB"),
	)
	token = os.Getenv("TOKEN")

}

func startDiscordSession() *discordgo.Session {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentGuildMembers | discordgo.IntentGuildPresences

	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord connection:", err)
	}
	return discord
}

func handleCommandLine() {
	kong.Parse(&cli)
	if cli.Version {
		ver, err := venova.GetVersionInfo("venova")
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Println(ver)
		os.Exit(0)
	}
}
