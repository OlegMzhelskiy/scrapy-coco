package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b TgBot) SendTextMessageWithRetry(Text string, ChatID int64, MessageID int, retryCount int) error {
	var err error

	msg := tgbotapi.NewMessage(ChatID, Text)
	if MessageID != 0 {
		msg.ReplyToMessageID = MessageID
	}

	for i := 0; i <= retryCount; i++ {
		_, err = b.Send(msg)
		if err == nil {
			return nil
		}
	}

	return err
}
