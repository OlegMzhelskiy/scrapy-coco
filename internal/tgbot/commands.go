package tgbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"scraper_nike/internal/log"
)

const (
	StartCommandText   = "Hi! 🙂\nI am a bot %s"
	HelpCommandText    = "I'm just forwarding messages to my master"
	UnknownCommandText = "Unknown command"
)

func (b TgBot) processCommand(message *tgbotapi.Message) {
	ChatID := message.Chat.ID
	MessageID := message.MessageID

	var err error

	command := message.Command()

	switch command {
	case "start":
		msg := tgbotapi.NewMessage(ChatID, fmt.Sprintf(StartCommandText, b.BotName))
		_, err = b.Send(msg)
	case "help":
		msg := tgbotapi.NewMessage(ChatID, HelpCommandText)
		_, err = b.Send(msg)
	default:
		msg := tgbotapi.NewMessage(ChatID, UnknownCommandText)
		msg.ReplyToMessageID = MessageID
		_, err = b.Send(msg)
	}

	if err != nil {
		log.GetLogger().Errorf("failed to send message from command '%s': %s", command, err)
	}
}
