package bot

import (
	"log"
	"os"

	"github.com/braydencw1/venova/db"

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
func getUserVoiceChannel(s *discordgo.Session, gId, uId string) string {
	guild, err := s.State.Guild(gId)
	if err != nil {
		log.Printf("cannot get guild state for voice channel: %s", err)
		return ""
	}
	for _, vs := range guild.VoiceStates {
		if vs.UserID == uId {
			return vs.ChannelID
		}
	}
	return ""
}

func GetEnvOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Printf("%s is empty or not defined. Defaulting to: %s", key, def)
		v = def
	}
	return v
}
