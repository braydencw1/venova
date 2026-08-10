package bot

import "slices"

func InitCommands() *CommandRegistry {
	cr := NewCommandRegistry()
	// Send birthday messages now
	cr.Register("bday", bdaySendCmd, 0, PermAdmin, "Usage !bday. Sends today's birthday messages immediately, if there are any. Admins only.")
	// Disconnect bot from VC
	cr.Register("dc", dcCmd, 0, PermAdmin, "Usage !dc <user> - Executes the disconnect user from VC command. Admins only.")
	// Set DND Date (Admin)
	cr.Register("dnd", playDndCmd, 1, PermAdmin, "Usage !dnd <args> - Executes the play dnd command. Updates next schedules play date. Admins only.")
	// Execute help command to display available
	// commands the caller is allowed to use.
	cr.Register("help", helpCmd, 0, PermEveryone, "Usage !help or !help <command>. - Displays all commands or command usage syntax.")
	// Execute minecraft command
	cr.Register("mc", mcCmd, 1, PermMcAdmin, "Usage !mc <args>. Executes Minecraft commands via RCON. Admins or Minecraft admins only.")
	// Restart Minecraft Server
	cr.Register("mcr", manageMinecraftCmd, 0, PermMcAdmin, "Usage !mcr <args>. Available: up, down, restart. Restarts the Minecraft server container. Admins or Minecraft admins.")
	// Play Audio Command
	cr.Register("play", playAudioCmd, 0, PermAdmin, "Usage !play. Allowing streaming of Audio to bots voice channel. Admins only.")
	// List joinable roles
	cr.Register("rlist", roleListCmd, 0, PermEveryone, "Usage !rlist. Lists available joinable roles via rjoin or rleave.")
	// Join role from list
	cr.Register("rjoin", roleJoinCmd, 1, PermEveryone, "Usage !rjoin <role>. Joins a joinable role.")
	// Leave role from list
	cr.Register("rleave", roleLeaveCmd, 1, PermEveryone, "Usage !rleave <role>. Leaves a joinable role.")
	// Roll dice
	cr.Register("roll", rollCmd, 1, PermEveryone, "Usage: !roll <dice> [adv/dis] or !roll stats. Examples: !roll 2d6+1d8+3, !roll d20 adv, !roll stats (4d6 drop lowest x6).")
	// Set a timer
	cr.Register("set", setTimerCmd, 1, PermAdmin, "Usage !set <duration> [@user]. Example: !set 1h30m. Sets a timer; the bot DMs the target when it's up. Survives restarts. Admins only.")
	// See when dnd is
	cr.Register("when", whenIsDndCmd, 0, PermEveryone, "Usage !when. Displays next DND play date if available to this discord server.")
	// Whitelist Minecraft
	cr.Register("whitelist", whitelistCmd, 1, PermMcAdmin, "Usage !whitelist <mcUserName>. Whitelists a Minecraft user to the server. Admins or Minecraft admins.")
	return cr
}

// ListCommandsFor returns the sorted names of commands the caller may use.
func (cr *CommandRegistry) ListCommandsFor(ctx CommandCtx) []string {
	keys := []string{}
	for name, cmd := range cr.commands {
		if cmd.perm.Allows(ctx) {
			keys = append(keys, name)
		}
	}
	slices.Sort(keys)
	return keys
}
