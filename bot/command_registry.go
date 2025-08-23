package bot

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandCtx struct {
	Session   *discordgo.Session
	Message   *discordgo.MessageCreate
	Args      []string
	IDChecker IdentityChecker
}

type Command struct {
	fn              func(c CommandCtx) error
	numRequiredArgs int
	help            string
}

type CommandRegistry struct {
	commands map[string]*Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*Command),
	}
}

func (c *CommandRegistry) Register(name string, command func(c CommandCtx) error, numArgs int, help string) {
	c.commands[name] = &Command{
		fn:              command,
		numRequiredArgs: numArgs,
		help:            help,
	}
}

func (c *CommandRegistry) HandleMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	parts := strings.SplitN(msg.Content, " ", 2)

	if len(parts) == 0 || len(parts[0]) == 0 {
		return
	}

	commandNameWithPrefix := strings.ToLower(parts[0])

	if commandNameWithPrefix[0] != '!' {
		return
	}
	commandName := commandNameWithPrefix[1:]

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
			_, err := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("the command %s has too few arguements.", commandName))
			if err != nil {
				log.Printf("err msgSend HandleMessage %s", err)
			}
			return
		}
	}
	ctx := CommandCtx{
		Session:   s,
		Message:   msg,
		Args:      args,
		IDChecker: GetIdentityChecker(),
	}

	go func() {
		if err := command.fn(ctx); err != nil {
			if err := ctx.Reply(fmt.Sprintf("error: %s", err)); err != nil {
				log.Printf("Handle MSGs err: %s", err)
			}
		}
	}()
}

func (c *CommandCtx) Reply(s string) error {
	_, err := c.Session.ChannelMessageSend(c.Message.ChannelID, s)
	return err
}

func (c *CommandCtx) DirectReply(s string) error {
	_, err := c.Session.ChannelMessageSend(c.Message.Author.ID, s)
	return err
}

func (c *CommandCtx) HasDiscordRole(givenRole string) (bool, error) {
	mem, err := c.Session.State.Member(c.Message.GuildID, c.Message.Author.ID)
	if err != nil {
		if replyErr := c.Reply(fmt.Sprintf("could not find member: %s", err)); replyErr != nil {
			return false, replyErr
		}
		return false, err
	}

	return slices.Contains(mem.Roles, givenRole), nil
}
