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
		insertRes, err := db.InsertPlayDate(t)
		if err != nil {
			log.Panic(err)
		}
		if insertRes {
			discord.ChannelMessageSend(tcDndGeneralId, "The Date has been updated.")
		}
	}
	if memberHasRole(member, dndRoleId) {
		if parts[0] == "!when" {
			res, err := db.GetLatestPlayDate()
			if err != nil {
				log.Printf("Retrieving last played date failed.")
			}
			fmtPlayDate := fmt.Sprint(res.Format("01-02-2006"))
			discord.ChannelMessageSend(tcDndGeneralId, fmtPlayDate)
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
