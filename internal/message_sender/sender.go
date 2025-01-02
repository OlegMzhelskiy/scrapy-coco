package message_sender

import (
	"errors"
	"fmt"
)

type transporter interface {
	SendTextMessageWithRetry(text string, chatID int64, amessageID int, retryCount int) error
}

type MessageSender struct {
	transporter transporter // tg
	chatIDs     []int64
	retryCount  int
}

func New(transporter transporter, chatIDs []int64, retryCount int) MessageSender {
	return MessageSender{
		transporter: transporter,
		chatIDs:     chatIDs,
		retryCount:  retryCount,
	}
}

// Send - messageID used for reply
func (m MessageSender) Send(message string, messageID int) error {
	var totalErr error

	for _, chatID := range m.chatIDs {
		if err := m.sendMessage(message, chatID, messageID); err != nil {
			totalErr = errors.Join(totalErr, err)
		}
	}

	if totalErr != nil {
		return fmt.Errorf("failed to send message: %w", totalErr)
	}

	return nil
}

func (m MessageSender) sendMessage(text string, chatID int64, messageID int) error {
	if err := m.transporter.SendTextMessageWithRetry(text, chatID, messageID, m.retryCount); err != nil {
		return fmt.Errorf("can't send message to chat id %v: %s", chatID, err)
	}

	return nil
}
