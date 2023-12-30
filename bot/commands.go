package bot

import (
	"fmt"
	"log"
	"strconv"
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

func memberHasDndRole(member *discordgo.Member) []string {
	var commonStuff []string
	res, err := db.GetDndRoles()
	if err != nil {
		log.Printf("Could not retrieve role Ids from DB")
	}
	for _, memberRoleId := range member.Roles {
		for _, arrayRoleId := range res {
			if memberRoleId == fmt.Sprintf("%v", arrayRoleId) {
				log.Printf("Your role is %v", memberRoleId)
				commonStuff = append(commonStuff, memberRoleId)
				break
			}
		}
	}
	log.Printf("Common stuff %v", commonStuff)
	return commonStuff
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
	if memberHasRole(member, dndRoleId) || memberHasRole(member, dndRoleIdH) {
		if parts[0] == "!when" {
			commonRoles := memberHasDndRole(member)
			if len(commonRoles) > 0 {
				for _, stringRole := range commonRoles {
					intRole, err := strconv.ParseInt(stringRole, 10, 64)
					if err != nil {
						log.Printf("Error convering role to int64, %v", err)
					}
					dateOfPlay, tcId, err := db.GetLatestPlayDate(intRole)
					if err != nil {
						log.Printf("Error parsing Latest playDate, %v", err)
					}
					log.Printf("dateofPlay %v", dateOfPlay)
					discord.ChannelMessageSend(fmt.Sprintf("%v", tcId), fmt.Sprint(dateOfPlay.Format("01-02-2006")))
				}
			} else {
				log.Printf("Not a part of a dnd role")
			}
		}
	}
	//			yn, tttt := memberHasDndRole(member, dndRoleId)
	//			if yn {
	//				log.Printf("logged")
	//				roleIdValue, err := strconv.ParseInt(tttt, 10, 64)
	//				if err != nil {
	//					fmt.Println("Error:", err)
	//					return
	//				}
	//				dateOfPlay, tcId, err := db.GetLatestPlayDate(roleIdValue)
	//				if err != nil {
	//					log.Printf("Retrieving last played date failed.")
	//				}
	//				fmtPlayDate := fmt.Sprint(dateOfPlay.Format("01-02-2006"))
	//				discord.ChannelMessageSend(fmt.Sprintf("%v", tcId), fmtPlayDate)
	//			}
	//		}
	//	}
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
