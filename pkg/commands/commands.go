package commands

import "discord-selfbot/pkg/selfbot"

// InitCommands creates a command list which is used by the selfbot.
// All commands should be added to the command list here
func InitCommands() selfbot.CommandList {
	list := selfbot.NewCommandList()

	list.AddCommand(selfbot.Command{
		Name:        "delete",
		Aliases:     []string{"d"},
		Run:         deleteCommand,
		Description: "Deletes the past n commands",
		Usage:       "delete [messageCount]",
		MinArgs:     0,
	})
	list.AddCommand(selfbot.Command{
		Name:        "ascii",
		Aliases:     nil,
		Run:         asciiCommand,
		Description: "Turns your text to ascii art and prints it to the chat",
		Usage:       "ascii <text>",
		MinArgs:     1,
	})
	list.AddCommand(selfbot.Command{
		Name:        "react",
		Aliases:     []string{"r"},
		Run:         reactCommand,
		Description: "Reacts custom text to the most recent message",
		Usage:       "react <text>",
		MinArgs:     1,
	})

	list.AddCommand(selfbot.Command{
		Name:        "help",
		Aliases:     []string{"h"},
		Run:         helpCommand,
		Description: "Prints help about all commands or about a specific command",
		Usage:       "help [command]",
		MinArgs:     0,
	})

	return list
}
