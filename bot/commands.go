package bot

import (
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"
	"time"
	"venova/db"

	"github.com/bwmarrin/discordgo"
)

var BotCommands = map[string]func(*discordgo.Session, *discordgo.MessageCreate, []string){
	//Minecraft Commands
	"mc": mcCmd,
	// Play Dnd Date (Admin)
	"play": playCmd,
	// Restart Minecraft Server
	"restart": restartMcCmd,
	// Join role from list
	"rjoin": roleJoinCmd,
	// Leave role from list
	"rleave": roleLeaveCmd,
	// List joinable roles
	"rlist": roleListCmd,
	// Set timer
	"set": setTimerCmd,
	//See when dnd is
	"when": whenIsDndCmd,
	//Whitelist Minecraft
	"whitelist": whitelistCmd,
}

func mcCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, parts []string) {
	if msg.Author.ID == blueId || msg.Author.ID == morthisId {
		res, err := minecraftCommand(parts[1])
		if err != nil {
			log.Printf("Err: %s", err)
			discord.ChannelMessageSend(msg.ChannelID, "Sorry you're not in a Minecraft Server Admin", nil)
		}
		discord.ChannelMessageSend(msg.ChannelID, res)
	}
}

func HandleCommands(discord *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID != discord.State.User.ID {
		log.Printf(msg.Author.Username + ": " + msg.Content)
	}

	parts := strings.SplitN(msg.Content, " ", 2)
	guild, err := discord.State.Guild(msg.GuildID)
	if err != nil {
		log.Printf("Failed to fetch message guild id (dm?): %v", err)
		return
	}

	member := getGuildMember(guild, msg.Author.ID)
	if member == nil {
		log.Printf("Failed to find message guild member: %v", msg.Author.ID)
		return
	}

	if fn, ok := BotCommands[fmt.Sprintf("!%s", parts[0])]; ok {
		fn(discord, msg, parts)
	} else {
		log.Printf("Invalid Command.")
		discord.ChannelMessageSend(msg.ChannelID, "Invalid command.", nil)
	}
}

func playCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, parts []string) {
	// Set dnd play date
	if msg.Author.ID == morthisId {
		layout := "01-02-2006"
		t, err := time.Parse(layout, parts[1])
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}
		currRoleId := getMemberDNDRole(msg.Member)
		if currRoleId == "" {
			log.Printf("Role not found.")
			discord.ChannelMessageSend(msg.ChannelID, "Your dnd role is not found in the db.")
		} else {
			err := db.InsertPlayDate(t, currRoleId)
			if err != nil {
				log.Panic(err)
			}
			discord.ChannelMessageSend(msg.ChannelID, "The Date has been updated.")

		}
	}
}
func restartMcCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, _ []string) {
	if memberHasRole(msg.Member, mcRoleId) {
		mcMsg, _ := discord.ChannelMessageSend(msg.ChannelID, "Restarting the minecraft server...")

		go func() {
			restartMinecraft()
			time.Sleep(time.Second * 5)
			discord.ChannelMessageEdit(msg.ChannelID, mcMsg.ID, "Minecraft server restarted!")
		}()
	}
}

func roleListCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, _ []string) {
	rolesString := strings.Join(slices.Collect(maps.Keys(joinableRolesMap)), ", ")
	discord.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Available roles: %s.\n Available commands: !rjoin & !rleave.", rolesString))
}

func roleLeaveCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, parts []string) {
	if roleID, exists := joinableRolesMap[parts[1]]; exists {
		err := discord.GuildMemberRoleRemove(msg.GuildID, msg.Author.ID, roleID)
		if err != nil {
			log.Printf("error removing role: %s", err)
		} else {
			log.Printf("Removed user with id: %s (%s) from %s role", msg.Author.ID, msg.Author.Username, roleID)
			discord.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("You've been removed from the group %s.", parts[1]))
		}
	}
}
func roleJoinCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, parts []string) {
	if roleID, exists := joinableRolesMap[parts[1]]; exists {
		err := discord.GuildMemberRoleAdd(msg.GuildID, msg.Author.ID, roleID)
		if err != nil {
			log.Printf("error adding role: %s", err)
		} else {
			log.Printf("Added user with id: %s (%s) to %s role", msg.Author.ID, msg.Author.Username, roleID)
			discord.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("You've been added to the group %s.", parts[1]))
		}
	}
}

func setTimerCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, parts []string) {
	if msg.Author.ID == morthisId || msg.Author.ID == bettyId {
		extraParts := strings.SplitN(parts[1], " ", 2)
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
		timerDestUserName, err := GetUsernameFromID(discord, extraParts[1])
		if err != nil {
			log.Printf("Could not retriever UserName from userID, %s", err)
		}
		log.Printf("Creating a timer for %s, for the length %s, destined for %s with userID %s", msg.Author.Username, extraParts[0], timerDestUserName, extraParts[1])
		if err != nil {
			log.Printf("%s", err)
		}
		errChan := make(chan error)
		go TimerCheckerRoutine(discord, timer, extraParts[1], errChan)
		err = <-errChan
		if err != nil {
			log.Printf("Error with the timer routine, %v", err)
		}
	}
}

func whenIsDndCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, parts []string) {
	now := time.Now()
	currRoleId := getMemberDNDRole(msg.Member)
	if currRoleId == "" {
		log.Printf("Could not find Dnd Role")
	}
	dateOfPlay, tcId, err := db.GetLatestPlayDate(currRoleId)
	if err != nil {
		log.Printf("Error parsing Latest playDate, %v", err)
	}
	fmtDate := fmt.Sprint(dateOfPlay.Format("01-02-2006"))
	s := fmt.Sprintf("%v", tcId)
	if dateOfPlay.Before(now) {
		discord.ChannelMessageSend(s, fmt.Sprintf("There is no date currently set. Your last session was: %s", fmtDate))
	} else {
		discord.ChannelMessageSend(s, fmtDate)
	}
}

func whitelistCmd(discord *discordgo.Session, msg *discordgo.MessageCreate, parts []string) {
	if memberHasRole(msg.Member, mcRoleId) {
		log.Printf("Whitelisting, %s ", parts[1])
		res, err := minecraftCommand(fmt.Sprintf("whitelist add %s", parts[1]))
		if err != nil {
			log.Printf("Err: %s", err)
		}
		discord.ChannelMessageSend(msg.ChannelID, res)
	}
}
