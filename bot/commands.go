package bot

import (
	"fmt"
	"log"
	"strings"
	"time"
	"venova/db"

	"github.com/bwmarrin/discordgo"
)

func memberHasRole(member *discordgo.Member, roleId string) bool {
	for _, memberRoleId := range member.Roles {
		if memberRoleId == roleId {
			return true
		}
	}
	return false
}

func getMemberDNDRole(member *discordgo.Member) string {
	res, err := db.GetDndRoles()
	if err != nil {
		log.Printf("Could not retrieve role Ids from DB")
	}
	for _, memberRoleId := range member.Roles {
		for _, arrayRoleId := range res {
			if memberRoleId == arrayRoleId {
				return memberRoleId
			}
		}
	}
	return ""
}
func getGuildMember(guild *discordgo.Guild, userId string) *discordgo.Member {
	var member *discordgo.Member
	for _, m := range guild.Members {
		if m.User.ID == userId {
			member = m
			break
		}
	}
	return member
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

	if parts[0] == "!play" && msg.Author.ID == morthisId {
		layout := "01-02-2006"
		t, err := time.Parse(layout, parts[1])
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}
		currRoleId := getMemberDNDRole(member)
		if currRoleId == "" {
			log.Printf("Role not found.")
		} else {
			err := db.InsertPlayDate(t, currRoleId)
			if err != nil {
				log.Panic(err)
			}
			discord.ChannelMessageSend(msg.ChannelID, "The Date has been updated.")

		}
	}
	if parts[0] == "!set" && (msg.Author.ID == morthisId || msg.Author.ID == bettyId) {
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
			log.Printf("Could not create timer %w", err)
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
	if parts[0] == "!when" {
		currRoleId := getMemberDNDRole(member)
		if currRoleId == "" {
			log.Printf("Could not find Dnd Role")
		} else {
			dateOfPlay, tcId, err := db.GetLatestPlayDate(currRoleId)
			if err != nil {
				log.Printf("Error parsing Latest playDate, %v", err)
			}
			discord.ChannelMessageSend(fmt.Sprintf("%v", tcId), fmt.Sprint(dateOfPlay.Format("01-02-2006")))
		}
	}
	if memberHasRole(member, mcRoleId) || memberHasRole(member, frostedRoleId) {
		if parts[0] == "!restart" {
			mcMsg, _ := discord.ChannelMessageSend(msg.ChannelID, "Restarting the minecraft server...")

			go func() {
				restartMinecraft()
				time.Sleep(time.Second * 5)
				discord.ChannelMessageEdit(msg.ChannelID, mcMsg.ID, "Minecraft server restarted!")
			}()
		} else if parts[0] == "!mc" && (msg.Author.ID == blueId || msg.Author.ID == morthisId) {
			res, err := minecraftCommand(parts[1])
			if err != nil {
				log.Printf("Err: %s", err)
			}
			discord.ChannelMessageSend(msg.ChannelID, res)
		} else if strings.HasPrefix(msg.Content, "!whitelist") {
			log.Printf("Whitelisting, %s ", parts[1])
			res, err := minecraftCommand(fmt.Sprintf("whitelist add %s", parts[1]))
			if err != nil {
				log.Printf("Err: %s", err)
			}
			discord.ChannelMessageSend(msg.ChannelID, res)
		}
	}
}
