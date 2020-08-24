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
	"sync"
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

	// for multiple configs. if there is only one config, use configs[0]
	configs := []selfbot.Config{{}}
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		// does not exist, ask for token and generate a new file
		reader := bufio.NewReader(os.Stdin)
		log.Infof("Could not find config file %s, creating a new one", configFileName)
		// TODO: add instructions for getting token
		fmt.Print("Enter your token: ")
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)

		// Create a new config
		configs[0] = selfbot.NewConfigDefault(token)

		// Save the new config file
		configFile, err := os.Create(configFileName)
		if err != nil {
			panic(err)
		}

		enc := json.NewEncoder(configFile)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		_ = enc.Encode(configs[0])
	} else {
		// Open and read config file
		configJson, err := ioutil.ReadFile(configFileName)
		if err != nil {
			log.Fatalln(err)
		}
		configs[0], err = selfbot.NewConfigFromJson(configJson)

		if err != nil {
			// try reading multiple configs
			configs, err = selfbot.NewConfigsFromJson(configJson)
			if err != nil {
				log.Fatalf("Error parsing config file %s: %s", configFileName, err)
			}
		}
	}

	// start a selfbot for every config
	commandList := commands.InitCommands()
	var selfbotList []selfbot.Selfbot
	for i, config := range configs {
		config := config
		go func() {
			bot, err := selfbot.NewSelfbot(config, commandList)
			if err != nil {
				log.Errorf("Error creating selfbot for config index %d: %s", i, err)
			}
			selfbotList = append(selfbotList, bot)
		}()
	}

	// Wait for ctrl c
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc


	var wg sync.WaitGroup
	// Cleanly close down all Discord sessions
	for _, bot := range selfbotList {
		bot := bot
		wg.Add(1)
		go func() {
			bot.Close()
			wg.Done()
		}()
	}

	wg.Wait()
}
