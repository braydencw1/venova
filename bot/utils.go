package bot

import (
	"log"

	"github.com/braydencw1/venova/db"

	"github.com/bwmarrin/discordgo"
)

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
