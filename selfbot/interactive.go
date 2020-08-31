package selfbot

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"sync"
)

// Interactive is an interactive session for selfbot users. While this session
// is active, the selfbot will not accept commands and instead input will be sent
// to this struct.
//
// To Read from an interactive session, you can use a bufio.Reader. Messages are
// separated by a null character ('\000')
type Interactive struct {
	channelID           string
	consumeMessages     bool
	readBuffer          []byte
	writeBuffer			bytes.Buffer
	bot                 *Selfbot
	closed              bool
	endMessage          string
	sentMessageIDs      []string
	sentMessageIdsMutex sync.Mutex
	reader              *bufio.Reader
}

// StartInteractive starts a new interactive session
func (bot *Selfbot) StartInteractive(channelID string, consumeMessages bool) *Interactive {
	session := &Interactive{
		bot:             bot,
		channelID:       channelID,
		consumeMessages: consumeMessages,
		endMessage:      bot.Config.Prefix,
	}

	session.reader = bufio.NewReader(session)

	go session.sendThread()

	bot.interactiveSession = session
	_ = bot.SendInfo(channelID, "Started a new interactive session", true)
	bot.Log.Debugln("Started a new interactive session")
	return session
}

func (i *Interactive) handleSendMessage(ev *discordgo.MessageCreate) (handleNext bool) {
	if i.closed {
		i.bot.interactiveSession = nil
		return true
	}

	// check if the bot sent the message
	i.sentMessageIdsMutex.Lock()
	for _, id := range i.sentMessageIDs {
		if id == ev.Message.ID {
			i.sentMessageIdsMutex.Unlock()
			return false
		}
	}
	i.sentMessageIdsMutex.Unlock()
	if ev.Content == i.endMessage {
		_ = i.Close()
		return false
	}
	if ev.ChannelID != i.channelID {
		return true
	}

	buf := append([]byte(ev.Content), '\n')
	i.readBuffer = append(i.readBuffer, buf...)
	i.bot.Log.Debugf("Received interactive command: %s", ev.Content)
	if i.consumeMessages {
		_ = i.bot.Session.ChannelMessageDelete(ev.Message.ChannelID, ev.Message.ID)
	}

	return false
}

// constantly reads from the buffer and sends messages
func (i *Interactive) sendThread() {
	for {
		b, err := i.writeBuffer.ReadBytes('\n')

		if i.closed {
			return
		}

		// if we hit an eof, break (no more newline)
		if err != nil {
			break
		}

		if len(b) == 0 {
			continue
		}

		i.sentMessageIdsMutex.Lock()
		msg, err := i.bot.Session.ChannelMessageSend(i.channelID, fmt.Sprintf("```%s```", string(b)))

		if err != nil {
			i.bot.Log.WithField("message", string(b)).Errorf("Error sending to channel from interactive session: %s", err)
			i.sentMessageIdsMutex.Unlock()
			continue
		}

		i.sentMessageIdsMutex.Unlock()
		i.sentMessageIDs = append(i.sentMessageIDs, msg.ID)
	}
}

// ReadString reads a single message from the session. If there are no messages
// waiting, it will wait until a message is sent. If the session is closed,
// eof will be true
func (i *Interactive) ReadString() (message string, eof bool) {
	for {
		msg, err := i.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return "", true
			}
			continue
		}
		message = msg[:len(msg)-1]
		break
	}

	return
}

func (i *Interactive) WriteString(message string) (eof bool) {
	if i.closed {
		return true
	}

	if _, err := i.Write(append([]byte(message), '\n')); err == io.EOF {
		return true
	}

	return false
}

// Read reads data from the read buffer which is updated every new message sent.
// This can be easily read from using a Reader
func (i *Interactive) Read(p []byte) (n int, err error) {
	if i.closed {
		return 0, io.EOF
	}
	if len(i.readBuffer) == 0 {
		return 0, nil
	}
	if len(p) == 0 {
		return 0, nil
	}

	n = copy(p, i.readBuffer)
	if len(i.readBuffer) == n {
		i.readBuffer = nil
	} else {
		i.readBuffer = i.readBuffer[n:]
	}
	return
}

// Write writes a response to the channel
func (i *Interactive) Write(p []byte) (n int, err error) {
	if i.closed {
		return 0, io.EOF
	}

	// keep an internal buffer so we only write on newline
	i.writeBuffer.Write(p)


	return len(p), err
}

func (i *Interactive) Close() error {
	i.bot.interactiveSession = nil
	i.closed = true
	_ = i.bot.SendInfo(i.channelID, "Ended interactive session", true)
	i.bot.Log.Debugln("Closed interactive session")
	return nil
}
