package selfbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Selfbot represents all information required to run a selfbot. Multiple selfbots
// can be created with this struct
type Selfbot struct {
	Session     *discordgo.Session
	User        *discordgo.User
	Config      Config
	CommandList CommandList
	Log         *log.Entry
}

func NewSelfbot(config Config, commands CommandList) (Selfbot, error) {
	// Create a new Session
	session, err := discordgo.New(config.Token)
	if err != nil {
		return Selfbot{}, err
	}

	// Open the Session
	if err = session.Open(); err != nil {
		return Selfbot{}, err
	}

	// get local User
	user, err := session.User("@me")
	if err != nil {
		return Selfbot{}, fmt.Errorf("error getting local User: %s", user)
	}

	logger := log.WithFields(log.Fields{"User": user.Username + "#" + user.Discriminator})

	logger.Infof("Started selfbot")

	// create the selfbot struct
	bot := Selfbot{
		Session:       session,
		User:          user,
		Config:        config,
		CommandList:   commands,
		Log:           logger,
	}

	// init the handlers
	if err := bot.initHandlers(); err != nil {
		return Selfbot{}, nil
	}

	return bot, nil
}

func (bot *Selfbot) Close() {
	bot.Log.Debugf("Closing Session")
	_ = bot.Session.Close()
}
