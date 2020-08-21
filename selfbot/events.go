package selfbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (bot *Selfbot) initHandlers() error {
	// add the handlers
	bot.session.AddHandler(bot.onMessageCreate)

	return nil
}

// Called when any user sends a message in any channel
func (bot *Selfbot) onMessageCreate(session *discordgo.Session, ev *discordgo.MessageCreate) {
	if ev.Author.ID == bot.user.ID {
		bot.onSendMessage(ev)
		return
	}
}

// Called whenever the local user sends a message
func (bot *Selfbot) onSendMessage(ev *discordgo.MessageCreate) {
	content := ev.Message.Content
	if strings.HasPrefix(content, bot.config.Prefix) {
		// remove the prefix from content
		content = content[len(bot.config.Prefix):]

		if content[0] == ' ' {
			bot.log.WithField("content", content).Debugf("Ignoring command because it had a space after the prefix")
			return
		}

		// get the command args
		parsed := strings.Split(content, " ")

		command := parsed[0]
		args := parsed[1:]

		bot.log.Debugf("Received command: %s, with args: %v", command, args)

		bot.log = bot.log.WithFields(logrus.Fields{
			"command": command,
			"args":    args,
		})

		// delete the command message
		if err := bot.session.ChannelMessageDelete(ev.ChannelID, ev.Message.ID); err != nil {
			bot.log.Errorf("Error deleting command message: %s", err)
		}

		// handle the command
		userError, discordError := bot.handleCommand(command, args, ev.Message)

		if discordError != nil {
			userError = fmt.Errorf("An internal error occurred: %s", discordError)
			bot.log.Errorf("Error handling command: %s", discordError)
		}

		if userError != nil {
			bot.log.Errorf("user error: %s", userError)
			message, err := bot.session.ChannelMessageSendComplex(ev.ChannelID, &discordgo.MessageSend{
				Embed: &discordgo.MessageEmbed{
					Title:       "Error",
					Description: userError.Error(),
					Color:       0xea5455,
				},
			})
			if err != nil {
				bot.log.Debugf("Error sending error to channel %s", err)
			}
			go func() {
				time.Sleep(5 * time.Second)
				if err := bot.session.ChannelMessageDelete(message.ChannelID, message.ID); err != nil {
					bot.log.Errorf("Error deleting error message")
				}
			}()
		}
	}
}
