package commands

import (
	"discord-selfbot/selfbot"
	"fmt"
	"github.com/bwmarrin/discordgo"
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
			Footer: &discordgo.MessageEmbedFooter{
				Text: "by draven#4562\ngithub.com/rmccrystal/discord-selfbot",
			},
			Fields: field,
		},
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}
