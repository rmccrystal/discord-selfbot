package commands

import (
	"discord-selfbot/pkg/selfbot"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func reactCommand(bot *selfbot.Selfbot, args []string, message *discordgo.Message) (userError, discordError error) {
	// get latest message
	messageHistory, err := bot.Session.ChannelMessages(message.ChannelID, 100, "", "", "")
	if err != nil {
		return nil, err
	}

	if len(messageHistory) == 0 {
		return fmt.Errorf("no messages found"), nil
	}

	reactMessage := messageHistory[0]

	reactText := strings.Join(args, " ")
	var reactEmojis []string
	var usedEmojis []string
	for _, char := range reactText {
		if emojiArr, ok := EmojiDict[char]; ok {
			// find an emoji that hasn't been used
			for _, emoji := range emojiArr {
				// check if the emoji was used already
				used := false
				for _, usedEmoji := range usedEmojis {
					if usedEmoji == emoji {
						used = true
						break
					}
				}
				// if it wasn't used append it to the emoji arrays
				if !used {
					reactEmojis = append(reactEmojis, emoji)
					usedEmojis = append(usedEmojis, emoji)
					break
				}
				// otherwise, continue the loop and check the other emojis
			}
		}
	}

	// react with the emojis
	for _, emoji := range reactEmojis {
		if err := bot.Session.MessageReactionAdd(reactMessage.ChannelID, reactMessage.ID, emoji); err != nil {
			bot.Log.Errorf("error adding reaction")
		}
	}

	return nil, nil
}

func asciiCommand(bot *selfbot.Selfbot, args []string, message *discordgo.Message) (userError, discordError error) {
	input := strings.Join(args, " ")

	// create a figure
	fig := figure.NewFigure(input, bot.Config.DefaultAsciiFont, false)

	// format the text as monospaced
	messageText := fmt.Sprintf("```%s```", fig.String())

	if _, err := bot.Session.ChannelMessageSend(message.ChannelID, messageText); err != nil {
		return nil, err
	}

	return nil, nil
}

// Deletes the past n messages from the User
func deleteCommand(bot *selfbot.Selfbot, args []string, message *discordgo.Message) (userError, discordError error) {
	// if there are no args count is 0, else parse args[0]
	var count int
	if len(args) == 0 {
		count = 1
	} else {
		var err error
		count, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("%s is not a valid number", args[0]), nil
		}
	}

	// delete the command if it's not already deleted
	_ = bot.Session.ChannelMessageDelete(message.ChannelID, message.ID)

	// delete `count` past messages
	deleteCount := 0
	lastDeletedMessageID := ""
	for deleteCount < count {
		messages, err := bot.Session.ChannelMessages(message.ChannelID, 100, lastDeletedMessageID, "", "")
		if err != nil {
			return nil, fmt.Errorf("error getting message history: %s", err)
		}

		for _, message := range messages {
			if deleteCount >= count {
				break
			}

			if message.Author.ID != bot.User.ID {
				continue
			}

			err = bot.Session.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				bot.Log.WithFields(logrus.Fields{
					"deleteCount":    deleteCount,
					"messageType":    message.Type,
					"messageContent": message.Content,
				}).Warnf("error deleting message %s: %s", message.Content, err)
				continue
			}
			lastDeletedMessageID = message.ID
			deleteCount += 1
		}
	}

	return nil, nil
}
