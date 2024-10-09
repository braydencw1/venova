package bot

import (
	"log"
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
