package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rmccrystal/discord-selfbot/commands"
	"github.com/rmccrystal/discord-selfbot/selfbot"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	debug          bool
	configFileName string
)

func main() {
	// Parse flags
	flag.BoolVar(&debug, "debug", false, "Enables debug mode")
	flag.StringVar(&configFileName, "config", "config.json", "Specifies the config file")

	flag.Parse()

	// Init the logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})
	log.SetLevel(log.InfoLevel)

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// Generate or get config
	var config selfbot.Config
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		// does not exist, ask for token and generate a new file
		reader := bufio.NewReader(os.Stdin)
		log.Infof("Could not find config file %s, creating a new one", configFileName)
		// TODO: add instructions for getting token
		fmt.Print("Enter your token: ")
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)

		// Create a new config
		config = selfbot.NewConfigDefault(token)

		// Save the new config file
		configFile, err := os.Create(configFileName)
		if err != nil {
			panic(err)
		}

		enc := json.NewEncoder(configFile)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		_ = enc.Encode(config)
	} else {
		// Open and read config file
		configJson, err := ioutil.ReadFile(configFileName)
		if err != nil {
			log.Fatalln(err)
		}
		config, err = selfbot.NewConfigFromJson(configJson)

		if err != nil {
			log.Fatalf("Error parsing config file %s: %s", configFileName, err)
		}
	}

	commandList := commands.InitCommands()
	bot, err := selfbot.NewSelfbot(config, commandList)

	if err != nil {
		log.Errorf("Error creating selfbot: %s", err)
		return
	}

	// Wait for ctrl c
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
