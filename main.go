package main

import (
	"fmt"

	"scraper_nike/internal/app"
	"scraper_nike/internal/config"
	"scraper_nike/internal/log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	log.InitLogger(cfg.LogLevel)
	l := log.GetLogger()

	a, err := app.NewApp("", cfg)
	if err != nil {
		l.Fatal(err)
	}

	if err = a.Run(); err != nil {
		l.Fatal(err)
	}
}
