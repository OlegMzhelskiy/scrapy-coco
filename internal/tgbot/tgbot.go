package tgbot

import (
	"fmt"

	"scraper_nike/internal/config"
	"scraper_nike/internal/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	*tgbotapi.BotAPI
	BotName     string
	botID       int64
	adminChatID int64
	retryCount  int
}

func NewBot(token string, cfg config.Config) (*TgBot, error) {
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
		BotAPI:      bot,
		BotName:     fmt.Sprintf("@%s", userBot.UserName),
		botID:       userBot.ID,
		adminChatID: cfg.AdminChatID,
		retryCount:  cfg.RetryCount,
	}, nil
}
