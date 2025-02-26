package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/braydencw1/venova"
	"github.com/braydencw1/venova/bot"
	"github.com/braydencw1/venova/db"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	showVersion := flag.Bool("version", false, "Display version information")
	showVersionShort := flag.Bool("v", false, "Display version information (short flag)")
	flag.Parse()
	if *showVersion || *showVersionShort {
		ver := venova.GetVersion("venova")
		// fmt.Println(ver)
		fmt.Printf("Name:\t\t%s\nVersion:\t%s\nGit revision:\t%s\nGit ref:\t%s\nGO version:\t%s\nBuilt:\t\t%s\nOS/Arch:\t%s/%s\n",
			ver.Name, ver.Version, ver.Revision, ver.Reference, ver.GoVers, ver.BuiltAt, ver.OS, ver.Arch)
		os.Exit(0)
	}
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
	log.Printf("Venova is online.")
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
	// Creates / registers all cmds
	cr := bot.InitCommands()
	discord.AddHandler(cr.HandleMessage)
	discord.AddHandler(bot.AddGriefer)

	go bot.BirthdateCheckRoutine(discord)
	//go bot.PlayDateCheckRoutine(discord)
	select {} // Block the main goroutine indefinitely
}
