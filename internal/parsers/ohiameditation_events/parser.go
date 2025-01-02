package ohiameditation_events

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scraper_nike/internal/log"
	"scraper_nike/internal/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Parser struct {
	urlEvents string
	eventName string
}

func New(urlEvents string, eventName string) Parser {
	return Parser{
		urlEvents: urlEvents,
		eventName: eventName,
	}
}

func (p Parser) Parse() ([]models.Event, error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithErrorf(log.GetLogger().Warnf),
		chromedp.WithLogf(log.GetLogger().Infof))
	defer cancel()

	events := make([]models.Event, 0)

	// exec task
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(p.urlEvents),
		chromedp.Sleep(15*time.Second), // await JavaScript
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		log.GetLogger().Errorf("failed to exec task: %s", err)
	}

	// parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.GetLogger().Errorf("failed to parse page: %s", err)
	}

	doc.Find(".bw-widget__day").
		Each(func(i int, day *goquery.Selection) {
			eventDay := day.Find(".bw-widget__date").Text()

			// find schedule
			day.Find(".bw-session").
				Each(func(i int, s *goquery.Selection) {
					title := s.Find(".bw-session__name").Text()
					if strings.Contains(strings.ToLower(title), p.eventName) {
						eventName := s.Find(".bw-session__name").Text()
						eventTime := s.Find(".bw-session__column").Text()
						instructor := s.Find(".bw-session__staff").Text()
						meta := s.Find(".bw-session__meta").Text()
						description := s.Find(".bw-session__description").Text()

						button := s.Find("[data-url]")
						dataUrl, _ := button.Attr("data-url")

						event := models.Event{
							Name:        prepareText(eventName),
							Description: prepareText(description),
							Date:        prepareText(eventDay),
							Time:        prepareText(eventTime),
							Instructor:  prepareText(instructor),
							Meta:        prepareText(meta),
							URL:         dataUrl,
						}

						events = append(events, event)

						param := fmt.Sprintf("Дата: %s, Время: %s, Класс: %s, Инструктор: %s, Дополнительно: %s, Описание: %s, URL: %s",
							eventDay, eventTime, title, instructor, meta, description, dataUrl)

						str := strings.ReplaceAll(strings.ReplaceAll(param, "\n", ""), "  ", "")

						log.GetLogger().Infof(" %s\n", str)
					}
				})
		})

	return events, nil
}

func prepareText(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), "  ", "")
}
