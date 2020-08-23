package main

import (
	"discord-selfbot/pkg/commands"
	"discord-selfbot/pkg/selfbot"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

const (
	ConfigFile = "config.json"
)

func main()  {
	// Init the logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})
	log.SetLevel(log.DebugLevel)

	// Open and read config file
	configFile, err := os.Open(ConfigFile)
	if err != nil {
		log.Fatalf("Error opening config file %s: %s", ConfigFile, err)
	}

	configJson, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatalf("Error reading config file %s: %s", ConfigFile, err)
	}

	config, err := selfbot.NewConfigFromJson(configJson)
	if err != nil {
		log.Fatalf("Error parsing config file %s: %s", ConfigFile, err)
	}

	commandList := commands.InitCommands()

	bot, err := selfbot.NewSelfbot(config, commandList)

	if err != nil {
		log.Errorf("Error creating selfbot: %s", bot)
		return
	}

	// Wait for ctrl c
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
