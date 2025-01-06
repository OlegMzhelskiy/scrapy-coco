package config

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Config struct {
	Token            string        `mapstructure:"tg_token"`
	DebugBot         bool          `mapstructure:"debug_bot"`
	LogLevel         string        `mapstructure:"log_level"`
	StoreType        string        `mapstructure:"store_type"`
	FileStoreName    string        `mapstructure:"file_store_name"`
	RetryCount       int           `mapstructure:"retry_count"`
	ScrapingInterval time.Duration `mapstructure:"scraping_interval"`
	ChatIDs          []int64       `mapstructure:"chat_ids"`
	AdminChatID      int64         `mapstructure:"admin_chat_id"`
	EventName        string        `mapstructure:"event_name"`
	URLEventSource   string        `mapstructure:"url_event_source"`
	DB               db
}

type db struct {
	DBHost     string `mapstructure:"db_host"`
	DBPort     string `mapstructure:"db_port"`
	DBName     string `mapstructure:"db_name"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
}

func LoadConfig(configPath string) (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("../../" + configPath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("")

	bindEnvs()

	viper.SetDefault("scraping_interval", "1m")
	viper.SetDefault("retry_count", 1)
	viper.SetDefault("log_level", "DEBUG")
	viper.SetDefault("debug_bot", true)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Error("config file not found, falling back to environment variables", err)
		} else {
			log.Error("error reading config file", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return config, nil
}

func bindEnvs() {
	err := viper.BindEnv("tg_token") // TG_TOKEN
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("admin_chat_id") // ADMIN_CHAT_ID
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("chat_ids") // CHAT_IDS
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("event_name") // EVENT_NAME
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("url_event_source") // URL_EVENT_SOURCE
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("db_host")
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("db_port")
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("db_name")
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("db_user")
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("db_password")
	if err != nil {
		log.Fatal(err)
	}

}
