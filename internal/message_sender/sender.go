package message_sender

import (
	"errors"
	"fmt"
)

type Transporter interface {
	SendTextMessageWithRetry(Text string, ChatID int64, MessageID int, retryCount int) error
}

type MessageSender struct {
	transporter Transporter //tg
	chatIDs     []int64
	retryCount  int
}

func New(transp Transporter, chatIDs []int64, retryCount int) MessageSender {
	return MessageSender{
		transporter: transp,
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

func (m MessageSender) sendMessage(Text string, ChatID int64, MessageID int) error {
	if err := m.transporter.SendTextMessageWithRetry(Text, ChatID, MessageID, m.retryCount); err != nil {
		return fmt.Errorf("can't send message to chat id %v: %s", ChatID, err)
	}

	return nil
}
