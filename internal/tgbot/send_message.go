package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b TgBot) SendTextMessageWithRetry(text string, chatID int64, messageID int, retryCount int) error {
	var err error

	msg := tgbotapi.NewMessage(chatID, text)
	if messageID != 0 {
		msg.ReplyToMessageID = messageID
	}

	for i := 0; i <= retryCount; i++ {
		_, err = b.Send(msg)
		if err == nil {
			return nil
		}
	}

	return err
}

func (b TgBot) SendMessageToAdmin(text string) error {
	textMsg := tgbotapi.NewMessage(b.adminChatID, text)

	_, err := b.Send(textMsg)

	return err
}
