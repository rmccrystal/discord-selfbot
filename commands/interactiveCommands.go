package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rmccrystal/discord-selfbot/selfbot"
	"io"
	"net"
)

func testInteractiveCommand(bot *selfbot.Selfbot, args []string, message *discordgo.Message) (userError, discordError error) {
	interactive := bot.StartInteractive(message.ChannelID, false)
	for {
		msg, eof := interactive.ReadString()
		if eof {
			break
		}

		if eof := interactive.WriteString(msg); eof {
			break
		}
	}
	return
}

func netcatCommand(bot *selfbot.Selfbot, args []string, message *discordgo.Message) (userError, discordError error) {
	conn, err := net.Dial("tcp", args[0])
	if err != nil {
		return err, nil
	}
	defer conn.Close()

	interactive := bot.StartInteractive(message.ChannelID, false)

	go io.Copy(interactive, conn)
	io.Copy(conn, interactive)

	//var wg sync.WaitGroup
	//go func() {
	
	//	wg.Add(1)
	//	if _, err := io.Copy(conn, interactive); err != nil {
	//		bot.Log.Errorln(err)
	//	}
	//	wg.Done()
	//}()
	//go func() {
	//	wg.Add(1)
	//	if _, err := io.Copy(interactive, conn); err != nil {
	//		bot.Log.Errorln(err)
	//	}
	//	wg.Done()
	//}()
	//
	//wg.Wait()

	return nil, nil
}
