package commands

import "github.com/rmccrystal/discord-selfbot/selfbot"

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
		Name:        "clear",
		Aliases:     []string{"c"},
		Run:         clearCommand,
		Description: "Sends x blank lines to the channel. Defaults to 60",
		Usage:       "clear [lines]",
	})
	list.AddCommand(selfbot.Command{
		Name: "clearPins",
		Run:  clearPinsCommand,
		Aliases: []string{"unpin"},
		Description: "Clears the oldest x pinned messages in the current channel or all pinned messages if no arguments are specified",
		Usage: "clearPins [messages]",
	})
	list.AddCommand(selfbot.Command{
		Name: "restorePins",
		Run: restorePinsCommand,
		Description: "Restores all pinned messages deleted by the clearPins command",
		Usage: "restorePins",
	})
	list.AddCommand(selfbot.Command{
		Name: "test",
		Run: testInteractiveCommand,
	})
	list.AddCommand(selfbot.Command{
		Name: "nc",
		Run: netcatCommand,
		MinArgs: 1,
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
