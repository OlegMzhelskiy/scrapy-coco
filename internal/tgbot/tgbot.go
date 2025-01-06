package tgbot

import (
	"context"
	"fmt"

	"scraper_nike/internal/config"
	"scraper_nike/internal/log"
	"scraper_nike/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type messageStore interface {
	GetMessageByID(ctx context.Context, msgID int) (models.TgMessage, error)
	SaveMessage(ctx context.Context, msg models.TgMessage) error
}

type TgBot struct {
	*tgbotapi.BotAPI
	BotName      string
	botID        int64
	adminChatID  int64
	retryCount   int
	messageStore messageStore
}

func NewBot(token string, cfg config.Config, store messageStore) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("can't creat bot: %w", err)
	}

	if err := tgbotapi.SetLogger(log.GetLogger()); err != nil {
		return nil, err
	}

	bot.Debug = cfg.DebugBot

	userBot, err := bot.GetMe()
	if err != nil {
		return nil, fmt.Errorf("can't get bot info: %w", err)
	}

	return &TgBot{
		BotAPI:       bot,
		BotName:      fmt.Sprintf("@%s", userBot.UserName),
		botID:        userBot.ID,
		adminChatID:  cfg.AdminChatID,
		retryCount:   cfg.RetryCount,
		messageStore: store,
	}, nil
}
