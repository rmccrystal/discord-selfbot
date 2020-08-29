package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rmccrystal/discord-selfbot/selfbot"
)

func helpCommand(bot *selfbot.Selfbot, args []string, message *discordgo.Message) (userError, discordError error) {
	var field []*discordgo.MessageEmbedField

	for _, command := range bot.CommandList.Commands {
		field = append(field, &discordgo.MessageEmbedField{
			Name:  bot.Config.Prefix + command.Name,
			Value: fmt.Sprintf("%s\nUsage: `%s%s`", command.Description, bot.Config.Prefix, command.Usage),
		})
	}

	message, err := bot.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:  "Help",
			Color:  0x4281A4,
			Fields: field,
			Footer: selfbot.Footer,
		},
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}
