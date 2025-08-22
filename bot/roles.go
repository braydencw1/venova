package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/braydencw1/venova/db"
)

func roleListCmd(ctx CommandCtx) error {
	roles, err := db.GetJoinableRoles(ctx.Message.GuildID)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("%s", err))
	}

	res := "Joinable roles include: \n"

	for _, v := range roles {
		res += v.Nickname + "\n"
	}
	return ctx.Reply(res)
}

func roleJoinCmd(ctx CommandCtx) error {
	// Change to compare against roles in the server
	args := ctx.Args
	msg := ctx.Message
	sess := ctx.Session

	var roleIDInt int64
	roles, err := db.GetJoinableRoles(ctx.Message.GuildID)
	if err != nil {
		return fmt.Errorf("issue retrieving joinable roles: %s", err)
	}
	for _, role := range roles {
		if strings.EqualFold(args[0], role.Nickname) {
			roleIDInt = role.RoleID
		}
	}
	roleID := strconv.Itoa(int(roleIDInt))

	hasRole, err := ctx.HasDiscordRole(roleID)

	if err != nil {
		return err
	}

	if hasRole {
		return nil
	}

	if err := sess.GuildMemberRoleAdd(msg.GuildID, msg.Author.ID, roleID); err != nil {
		log.Printf("error adding role: %s", err)
		return fmt.Errorf("error adding user to role: %w", err)
	}

	log.Printf("Added user with id: %s (%s) to %s role", msg.Author.ID, msg.Author.Username, roleID)
	return ctx.Reply(fmt.Sprintf("You've been added to the group %s.", args[0]))
}

func roleLeaveCmd(ctx CommandCtx) error {
	args := ctx.Args
	msg := ctx.Message

	var roleIDInt int64
	roles, err := db.GetJoinableRoles(ctx.Message.GuildID)
	if err != nil {
		return fmt.Errorf("issue retrieving joinable roles: %s", err)
	}
	for _, role := range roles {
		if strings.EqualFold(args[0], role.Nickname) {
			roleIDInt = role.RoleID
		}
	}
	roleID := fmt.Sprintf("%d", roleIDInt)

	hasRole, err := ctx.HasDiscordRole(roleID)

	if err != nil {
		return err
	}

	if !hasRole {
		return nil
	}

	if err := ctx.Session.GuildMemberRoleRemove(msg.GuildID, msg.Author.ID, roleID); err != nil {
		return fmt.Errorf("could not remove role: %w", err)
	}

	log.Printf("Removed user with id: %s (%s) from %s role", msg.Author.ID, msg.Author.Username, roleID)
	return ctx.Reply(fmt.Sprintf("You've been removed from the group %s.", args[0]))
}
