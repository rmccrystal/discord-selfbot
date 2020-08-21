package selfbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Selfbot represents all information requred to run a selfbot. Multiple selfbots
// can be created with this struct
type Selfbot struct {
	session *discordgo.Session
	user    *discordgo.User
	config  Config
	log     *log.Entry
}

func NewSelfbot(config Config) (Selfbot, error) {
	// Create a new session
	session, err := discordgo.New(config.Token)
	if err != nil {
		return Selfbot{}, err
	}

	// Open the session
	if err = session.Open(); err != nil {
		return Selfbot{}, err
	}

	// get local user
	user, err := session.User("@me")
	if err != nil {
		return Selfbot{}, fmt.Errorf("error getting local user: %s", user)
	}

	logger := log.WithFields(log.Fields{"user": user.Username + "#" + user.Discriminator})

	logger.Infof("Started selfbot")

	// create the selfbot struct
	bot := Selfbot{
		session: session,
		user:    user,
		config:  config,
		log:     logger,
	}

	// init the handlers
	if err := bot.initHandlers(); err != nil {
		return Selfbot{}, nil
	}

	return bot, nil
}

func (bot *Selfbot) Close() {
	bot.log.Debugf("Closing session")
	_ = bot.session.Close()
}
