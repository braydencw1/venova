package bot

func InitCommands() *CommandRegistry {
	cr := NewCommandRegistry()
	// Disconnect bot from VC
	cr.Register("dc", dcCmd, 0, "Usage !dc <user> - Executes the disconnect user from VC command. Admins only.")
	// Set DND Date (Admin)
	cr.Register("dnd", playDndCmd, 1, "Usage !dnd <args> - Executes the play dnd command. Updates next schedules play date. Admins only.")
	// Execute help command to display available
	// Commands. Currently list all, even if
	// unavailable to the user via roles etc.
	cr.Register("help", helpCmd, 0, "Usage !help or !help <command>. - Displays all comamnds or command usage syntax.")
	// Execute minecraft command
	cr.Register("mc", mcCmd, 1, "Usage !mc <args>. Executes Minecraft commands via RCON. Admins or Minecraft admins only.")
	// Restart Minecraft Server
	cr.Register("mcr", manageMinecraftCmd, 0, "Usage !mrc <args>. Available: up, down, restart. Restarts the Minecraft server container. Admins or Minecraft admins.")
	// Play Audio Command
	cr.Register("play", playAudioCmd, 0, "Usage !play. Allowing streaming of Audio to bots voice channel. Admins only.")
	// List joinable roles
	cr.Register("rlist", roleListCmd, 0, "Usage !rlist. Lists available joinable roles via rjoin or rleave.")
	// Join role from list
	cr.Register("rjoin", roleJoinCmd, 1, "Usage !rjoin <role>. Joins a joinable role.")
	// Leave role from list
	cr.Register("rleave", roleLeaveCmd, 1, "Usage !rleave <role>. Leaves a joinable role.")
	// Roll a die or dice
	cr.Register("roll", rollCmd, 1, "Usage: !roll [dice] [adv/dis]. Example: !roll 3d6+3 or !roll d20 adv")
	// Set timer
	cr.Register("set", setTimerCmd, 1, "Usage !set <00h-00m>. Set a timer in which the bot will DM you when the timer is up.")
	// See when dnd is
	cr.Register("when", whenIsDndCmd, 0, "Usage !when. Displays next DND play date if available to this discord server.")
	// Whitelist Minecraft
	cr.Register("whitelist", whitelistCmd, 1, "Usage !whitelist <mcUserName>. Whitelists a Minecraft user to the server. Admins or Minecraft admins.")
	return cr
}

func (cr *CommandRegistry) ListCommands() []string {
	keys := []string{}
	for k := range cr.commands {
		keys = append(keys, k)
	}
	return keys
}
