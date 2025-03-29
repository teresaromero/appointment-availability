package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	HLA struct {
		BaseURL            string `env:"HLA_BASE_URL"`
		Username           string `env:"HLA_USERNAME"`
		Password           string `env:"HLA_PASSWORD"`
		HealthCentreIDList []int  `env:"HLA_HEALTH_CENTRE_ID_LIST"`
		SpecialtyIDList    []int  `env:"HLA_SPECIALTY_ID_LIST"`
	}
	Telegram struct {
		Apikey string `env:"TG_BOT_APIKEY"`
		ChatID int64  `env:"TG_BOT_MASTERID"`
	}
	HC struct {
		URL                string   `env:"HC_URL"`
		HealthCentreIDList []string `env:"HC_HEALTH_CENTRE_ID_LIST"`
		SpecialtyIDList    []string `env:"HC_SPECIALTY_ID_LIST"`
	}
}

func Load() (*Config, error) {
	c := &Config{}
	if err := godotenv.Load(); err != nil {
		log.Default().Printf("Error loading .env file: %v", err)
	}

	if err := env.Parse(c); err != nil {
		return nil, fmt.Errorf("Error parsing config: %v", err)
	}
	return c, nil
}
