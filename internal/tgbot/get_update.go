package tgbot

import (
	"context"
	"fmt"
	"time"

	"scraper_nike/internal/log"
	"scraper_nike/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const megTemplate = `[%s]%s
From: %s %s (%d)
UserName: %s
Message: %s`

func (b TgBot) GetUpdateMessage() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	for update := range updates {
		b.processMessage(context.Background(), update)
	}
}

func (b TgBot) processMessage(ctx context.Context, update tgbotapi.Update) {
	if update.Message != nil {
		if !b.IsMessageToMe(*update.Message) && update.Message.Chat.Type != "private" {
			return
		}

		message := update.Message
		msg := messageFromTg(message)

		log.GetLogger().Infof("username: %s ID: %v text: %s", msg.UserName, msg.FromID, msg.Text)

		var warning string

		if err := b.messageStore.SaveMessage(ctx, msg); err != nil {
			warning = "‼️"
			log.GetLogger().Errorf("failed to save message to db: %s", err)
		}

		if message.IsCommand() {
			b.processCommand(message)

			return
		}

		text := fmt.Sprintf(megTemplate, msg.Date, warning, msg.FirstName,
			msg.LastName, msg.FromID, msg.UserName, msg.Text)

		textMsg := tgbotapi.NewMessage(b.adminChatID, text)

		sentMessage, err := b.Send(textMsg)
		if err != nil {
			log.GetLogger().Errorf("failed to send message: %s", err)

			return
		}

		log.GetLogger().Infof("returned message: %s", sentMessage)
	}
}

func messageFromTg(message *tgbotapi.Message) models.TgMessage {
	msg := models.TgMessage{
		ID:   message.MessageID,
		Text: message.Text,
		Date: time.Unix(int64(message.Date), 0).UTC(),
	}

	if message.From != nil {
		msg.UserName = message.From.UserName
		msg.FirstName = message.From.FirstName
		msg.LastName = message.From.LastName
		msg.FromID = int(message.From.ID)
	}

	if message.ReplyToMessage != nil {
		msg.ReplyToMessageID = message.ReplyToMessage.MessageID
	}

	if message.Chat != nil {
		msg.ChatID = int(message.Chat.ID)
	}

	return msg
}
