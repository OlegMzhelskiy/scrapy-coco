package app

import (
	"fmt"
	"time"

	"scraper_nike/internal/config"
	"scraper_nike/internal/log"
	"scraper_nike/internal/message_sender"
	events "scraper_nike/internal/parsers/ohiameditation_events"
	"scraper_nike/internal/store"
	"scraper_nike/internal/tgbot"
	"scraper_nike/internal/worker"
)

type ParseWorker interface {
	Run() error
}

type MessageSender interface {
	SendMessageToAdmin(text string) error
}

type App struct {
	done             chan struct{}
	tgBot            *tgbot.TgBot
	retryCount       int
	UrlEvents        string
	scrapingInterval time.Duration
	ChatIDs          []int64
	AdminChatID      int64
	EventName        string
	parseWorker      ParseWorker
	messageSender    MessageSender
}

func NewApp(token string, cfg config.Config) (*App, error) {
	if token == "" {
		token = cfg.Token
	}

	if token == "" {
		return nil, fmt.Errorf("token shouldn't be empty")
	}

	if cfg.AdminChatID == 0 {
		return nil, fmt.Errorf("admin_chat_id shouldn't be empty")
	}

	bot, err := tgbot.NewBot(token, cfg)
	if err != nil {
		return nil, err
	}

	doneCh := make(chan struct{})

	a := &App{
		done:             doneCh,
		tgBot:            bot,
		retryCount:       cfg.RetryCount,
		scrapingInterval: cfg.ScrapingInterval,
		ChatIDs:          cfg.ChatIDs,
		AdminChatID:      cfg.AdminChatID,
		EventName:        cfg.EventName,
		UrlEvents:        cfg.URLEventSource,
		parseWorker: worker.NewWorker(
			doneCh,
			cfg.ScrapingInterval,
			message_sender.New(bot, cfg.ChatIDs, cfg.RetryCount),
			store.NewMemoryStore(),
			events.New(cfg.URLEventSource, cfg.EventName)),
		messageSender: bot,
	}

	return a, nil
}

func (a App) Run() error {
	log.GetLogger().Infof("start bot %s\n", a.tgBot.BotName)

	if err := a.messageSender.SendMessageToAdmin("bot is successfully started"); err != nil {
		log.GetLogger().Errorf("failed to send message to admin: %s", err)

		return err
	}

	go a.tgBot.GetUpdateMessage()

	if err := a.parseWorker.Run(); err != nil {
		return err
	}

	return nil
}
