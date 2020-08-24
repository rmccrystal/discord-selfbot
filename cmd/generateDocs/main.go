// This file generates a list of commands for the readme and generates the example config file
// using the hard coded default configuration
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rmccrystal/discord-selfbot/commands"
	"github.com/rmccrystal/discord-selfbot/selfbot"
	"io/ioutil"
	"os"
	"strings"
)

const ConfigOutput = "../../config.example.json"

var CommandListOutput = os.Stdout

func main() {
	commandList := commands.InitCommands()
	commandMD := generateCommandMD(commandList)
	_, _ = CommandListOutput.Write([]byte(commandMD))

	configBytes := generateExampleConfig()
	if err := ioutil.WriteFile(ConfigOutput, configBytes, 0644); err != nil {
		panic(err)
	}
}

// generateCommandMD generates an md file string containing a formatted list
// of all commands the selfbot has
func generateCommandMD(commandList selfbot.CommandList) string {
	var md string
	for _, command := range commandList.Commands {
		md += fmt.Sprintf("### `%s`\n\n"+
			"Description: `%s`\n\n"+
			"Usage: `%s`\n\n", command.Name, command.Description, command.Usage)

		if len(command.Aliases) > 0 {
			md += fmt.Sprintf("Aliases: `%s`\n\n", strings.Join(command.Aliases, ", "))
		}
	}

	return md
}

func generateExampleConfig() []byte {
	config := selfbot.NewConfigDefault("Your token here")

	buf := bytes.Buffer{}

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(config); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
