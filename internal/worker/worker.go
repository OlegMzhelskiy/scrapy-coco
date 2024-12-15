package worker

import (
	"errors"
	"time"

	"scraper_nike/internal/log"
	"scraper_nike/internal/models"
	"scraper_nike/internal/store"
)

type MessageSender interface {
	Send(message string, messageID int) error
}

type Storer interface {
	Save(e models.Event) error
	Get(key string) (models.Event, error)
}

type PageParser interface {
	Parse() ([]models.Event, error)
}

type ParseWorker struct {
	done             chan struct{}
	scrapingInterval time.Duration
	MessageSender    MessageSender
	EventsStore      Storer
	PageParser       PageParser
}

func NewWorker(done chan struct{}, scrapingInterval time.Duration, msgSender MessageSender, eventsStore Storer, pageParser PageParser) ParseWorker {
	return ParseWorker{
		done:             done,
		scrapingInterval: scrapingInterval,
		MessageSender:    msgSender,
		EventsStore:      eventsStore,
		PageParser:       pageParser,
	}
}

func (w ParseWorker) Run() error {
	ticker := time.NewTicker(w.scrapingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.startParse(); err != nil {
				log.GetLogger().Errorf("failed to find event: %s", err)
			}
		case <-w.done:
			log.GetLogger().Info("Shutting down gracefully")

			return nil
		}
	}
}

// StartParse find and send events
func (w ParseWorker) startParse() error {
	events, err := w.PageParser.Parse()
	if err != nil {
		log.GetLogger().Errorf("failed to parse html doc: %s", err)
	}

	for _, e := range events {
		_, err = w.EventsStore.Get(e.Key())
		if err == nil {
			continue // skip events already sent
		}

		if !errors.Is(err, store.ErrNotFound) {
			log.GetLogger().Errorf("failed to get event from store: %s", err)
		}

		if err := w.MessageSender.Send(e.String(), 0); err != nil {
			log.GetLogger().Errorf("failed to send message: %s", err)

			continue
		}

		if err = w.EventsStore.Save(e); err != nil {
			log.GetLogger().Errorf("failed to save event: %s", err)
		}

	}

	return err
}
