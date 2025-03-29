package main

import (
	"appointment-availability/internal/bot"
	"appointment-availability/internal/config"
	hcservice "appointment-availability/internal/services/hc"
	hlaservice "appointment-availability/internal/services/hla"
	"context"
	"log"
	"net/http"
	"time"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tgBot := bot.New(cfg.Telegram.Apikey, cfg.Telegram.ChatID)
	defer tgBot.Close(ctx)

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	// HLA services
	hla := hlaservice.New(
		client,
		cfg.HLA.BaseURL,
		cfg.HLA.Username,
		cfg.HLA.Password,
		tgBot.SendNotification,
	)

	if err := hla.Run(ctx, cfg.HLA.HealthCentreIDList, cfg.HLA.SpecialtyIDList); err != nil {
		log.Default().Printf("Error running HLA service: %v", err)
	}

	hc := hcservice.New(
		cfg.HC.URL,
		tgBot.SendNotification,
	)
	if err := hc.Run(ctx); err != nil {
		log.Default().Printf("Error running HC service: %v", err)
	}
}
