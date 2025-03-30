package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type HLAConfig struct {
	BaseURL            string `env:"HLA_BASE_URL"`
	Username           string `env:"HLA_USERNAME"`
	Password           string `env:"HLA_PASSWORD"`
	HealthCentreIDList []int  `env:"HLA_HEALTH_CENTRE_ID_LIST"`
	SpecialtyIDList    []int  `env:"HLA_SPECIALTY_ID_LIST"`
}

type HCConfig struct {
	BaseURL            string   `env:"HC_URL"`
	HealthCentreIDList []string `env:"HC_HEALTH_CENTRE_ID_LIST"`
	SpecialtyIDList    []string `env:"HC_SPECIALTY_ID_LIST"`
}

type TelegramConfig struct {
	Apikey string `env:"TG_BOT_APIKEY"`
	ChatID int64  `env:"TG_BOT_MASTERID"`
}

type Config struct {
	HLA      HLAConfig
	HC       HCConfig
	Telegram TelegramConfig
}

func Load() (*Config, error) {
	c := &Config{}
	// Load environment variables from .env file if it exists
	// This is useful for local development
	if err := godotenv.Load(); err != nil {
		log.Default().Printf("Error loading .env file: %v", err)
	}

	hla := &HLAConfig{}
	if err := env.Parse(hla); err != nil {
		return nil, fmt.Errorf("Error parsing hla: %v", err)
	}
	hc := &HCConfig{}
	if err := env.Parse(hc); err != nil {
		return nil, fmt.Errorf("Error parsing hc: %v", err)
	}
	telegram := &TelegramConfig{}
	if err := env.Parse(telegram); err != nil {
		return nil, fmt.Errorf("Error parsing telegram: %v", err)
	}
	c.HLA = *hla
	c.HC = *hc
	c.Telegram = *telegram
	return c, nil
}
