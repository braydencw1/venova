package bot

func InitCommands() *CommandRegistry {
	cr := NewCommandRegistry()
	// Disconnect bot from VC
	cr.Register("dc", dcCmd, 0)
	// Set DND Date (Admin)
	cr.Register("dnd", playDndCmd, 1)
	// Execute minecraft command
	cr.Register("mc", mcCmd, 1)
	// Restart Minecraft Server
	cr.Register("mcr", manageMinecraftCmd, 0)
	// Play Audio Command
	cr.Register("play", playAudioCmd, 0)
	// List joinable roles
	cr.Register("rlist", roleListCmd, 0)
	// Join role from list
	cr.Register("rjoin", roleJoinCmd, 1)
	// Leave role from list
	cr.Register("rleave", roleLeaveCmd, 1)
	// Set timer
	cr.Register("set", setTimerCmd, 1)
	// See when dnd is
	cr.Register("when", whenIsDndCmd, 0)
	// Whitelist Minecraft
	cr.Register("whitelist", whitelistCmd, 1)
	return cr
}
