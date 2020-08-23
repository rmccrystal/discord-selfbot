package main

import (
	"discord-selfbot/pkg/commands"
	"discord-selfbot/pkg/selfbot"
	"flag"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

var (
	debug      bool
	configFile string
)

func main() {
	// Parse flags
	flag.BoolVar(&debug, "debug", false, "Enables debug mode")
	flag.StringVar(&configFile, "config", "config.json", "Specifies the config file")

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

	// Open and read config file
	configJson, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config file %s: %s", configFile, err)
	}

	config, err := selfbot.NewConfigFromJson(configJson)
	if err != nil {
		log.Fatalf("Error parsing config file %s: %s", configFile, err)
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
