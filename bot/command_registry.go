package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandCtx struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Args    []string
}

type Command struct {
	fn              func(c CommandCtx) error
	numRequiredArgs int
}

type CommandRegistry struct {
	commands map[string]*Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*Command),
	}
}

func (c *CommandRegistry) Register(name string, command func(c CommandCtx) error, numArgs int) {
	c.commands[name] = &Command{
		fn:              command,
		numRequiredArgs: numArgs,
	}
}
func (c *CommandRegistry) HandleMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID != s.State.User.ID {
		log.Printf(msg.Author.Username + ": " + msg.Content)
	}

	parts := strings.SplitN(msg.Content, " ", 2)
	commandNameWithPrefix := strings.ToLower(parts[0])

	if commandNameWithPrefix[0] != '!' {
		return
	}
	commandName := commandNameWithPrefix[1:]
	// Finding command name
	command := c.commands[commandName]

	if command == nil {
		return
	}

	args := []string{}
	if len(parts) > 1 {
		args = strings.Split(parts[1], " ")
	}
	if command.numRequiredArgs > 0 {
		if len(args) < command.numRequiredArgs {
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("the command %s has too few arguements.", commandName))
			return
		}
	}
	ctx := CommandCtx{
		Session: s,
		Message: msg,
		Args:    args,
	}
	// Found mcCmd and executing
	err := command.fn(ctx)
	if err != nil {
		log.Printf("%s", err)
	}

}

func (c *CommandCtx) Reply(s string) error {
	_, err := c.Session.ChannelMessageSend(c.Message.ChannelID, s)
	return err
}
