package selfbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandList struct {
	Commands []Command
}

func NewCommandList() CommandList {
	return CommandList{}
}

func (list *CommandList) AddCommand(command Command) {
	// check if there are any Commands with the same name or alias
	for _, cmd := range list.Commands {
		// check aliases
		found := false
		var conflictingCommand Command

		if strings.ToLower(cmd.Name) == strings.ToLower(command.Name) {
			found = true
			conflictingCommand = cmd
		} else {
			for i := range cmd.Aliases {
				if found {
					break
				}
				for j := range command.Aliases {
					if strings.ToLower(cmd.Aliases[i]) == strings.ToLower(command.Aliases[j]) {
						found = true
						conflictingCommand = cmd
						break
					}
				}
			}

		}

		if found {
			panic(fmt.Sprintf("Error adding command %s because it or one of its aliases conflicts with commannd %s", command.Name, conflictingCommand.Name))
		}
	}

	if command.Usage == "" {
		command.Usage = command.Name
	}

	// add command to the list
	list.Commands = append(list.Commands, command)
}

// FindCommand finds a command in the commandlist using the name or alias.
// If no commands with the specified command or alias exist, it will return nil.
func (list CommandList) FindCommand(command string) *Command {
	for _, cmd := range list.Commands {
		// check of command matches
		if strings.ToLower(cmd.Name) == strings.ToLower(command) {
			return &cmd
		}

		// check if aliases match
		for _, alias := range cmd.Aliases {
			if strings.ToLower(alias) == strings.ToLower(command) {
				return &cmd
			}
		}

		// check next command
	}

	return nil
}

// Run finds and runs a command based on the command and arguments
func (list CommandList) Run(bot *Selfbot, command string, args []string, message *discordgo.Message) (userError, err error) {
	cmd := list.FindCommand(command)

	if cmd == nil {
		// if there no Commands, return a user error
		return fmt.Errorf("Unknown command: %s", command), nil
	}

	if len(args) < cmd.MinArgs {
		return fmt.Errorf("%s requires %d args while only %d was specified", command, cmd.MinArgs, len(args)), nil
	}

	return cmd.Run(bot, args, message)
}

type Command struct {
	Name        string
	Aliases     []string
	Run         func(bot *Selfbot, args []string, message *discordgo.Message) (userError, err error)
	Description string
	Usage		string
	MinArgs     int
}
