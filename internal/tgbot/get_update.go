package tgbot

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"scraper_nike/internal/log"
)

const megTemplate = `[%s]
From: %s %s (%d)
UserName: %s
Message: %s`

func (b TgBot) GetUpdateMessage() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	for update := range updates {
		b.processMessage(update)
	}
}

func (b TgBot) processMessage(update tgbotapi.Update) {
	if update.Message != nil {
		if !b.IsMessageToMe(*update.Message) && update.Message.Chat.Type != "private" {
			return
		}

		message := update.Message
		mesText := message.Text

		log.GetLogger().Infof("username: %s ID: %v text: %s", message.From.UserName, message.From.ID, mesText)

		if message.IsCommand() {
			b.processCommand(message)

			return
		}

		text := fmt.Sprintf(megTemplate, time.Unix(int64(message.Date), 0).UTC(), message.From.FirstName,
			message.From.LastName, message.From.ID, message.From.UserName, message.Text)

		textMsg := tgbotapi.NewMessage(b.adminChatID, text)

		msg, err := b.Send(textMsg)
		if err != nil {
			log.GetLogger().Errorf("failed to send message: %s", err)

			return
		}

		log.GetLogger().Infof("returned message: %s", msg)
	}
}
