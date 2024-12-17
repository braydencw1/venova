package bot

import (
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"
)

func roleListCmd(ctx CommandCtx) error {
	rolesString := strings.Join(slices.Collect(maps.Keys(joinableRolesMap)), ", ")
	return ctx.Reply(fmt.Sprintf("Available roles: %s.\n Available commands: !rjoin & !rleave.", rolesString))
}

func roleJoinCmd(ctx CommandCtx) error {
	args := ctx.Args
	msg := ctx.Message
	sess := ctx.Session
	if roleID, exists := joinableRolesMap[args[0]]; exists {
		err := sess.GuildMemberRoleAdd(msg.GuildID, msg.Author.ID, roleID)
		if err != nil {
			log.Printf("error adding role: %s", err)
			return fmt.Errorf("error adding user to role: %w", err)
		}
		log.Printf("Added user with id: %s (%s) to %s role", msg.Author.ID, msg.Author.Username, roleID)
		return ctx.Reply(fmt.Sprintf("You've been added to the group %s.", args[0]))
	}
	return nil
}

func roleLeaveCmd(ctx CommandCtx) error {
	args := ctx.Args
	msg := ctx.Message
	if roleID, exists := joinableRolesMap[args[0]]; exists {
		err := ctx.Session.GuildMemberRoleRemove(msg.GuildID, msg.Author.ID, roleID)
		if err != nil {
			return fmt.Errorf("error removing role: %w", err)
		}
		log.Printf("Removed user with id: %s (%s) from %s role", msg.Author.ID, msg.Author.Username, roleID)
		return ctx.Reply(fmt.Sprintf("You've been removed from the group %s.", args[0]))
	}
	return nil
}
