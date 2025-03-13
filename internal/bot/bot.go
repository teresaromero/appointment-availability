package bot

import (
	"context"
	"log"
	"os"
	"strconv"

	tgBot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AppointmentBot struct {
	bot      *tgBot.Bot
	masterID int64
}

func mustGetEnvInt64(name string) int64 {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("env variable %s is not set", name)
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("error parsing int64: %v", err)
	}
	return i
}

func New() *AppointmentBot {
	token := os.Getenv("TG_BOT_APIKEY")
	if token == "" {
		log.Fatal("TG_BOT_APIKEY is not set")
	}
	masterID := mustGetEnvInt64("TG_BOT_MASTERID")

	middleware := func(next tgBot.HandlerFunc) tgBot.HandlerFunc {
		return func(ctx context.Context, b *tgBot.Bot, update *models.Update) {
			if update.Message.From.ID != masterID ||
				update.Message.Chat.Type != "private" ||
				update.Message.Chat.ID != masterID {
				if _, err := b.SendMessage(ctx, &tgBot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Sorry this is private bot",
				}); err != nil {
					log.Printf("error sending message: %v", err)
				}
				return
			}
		}
	}

	opts := []tgBot.Option{
		tgBot.WithMiddlewares(middleware),
	}

	b, err := tgBot.New(token, opts...)
	if err != nil {
		log.Fatal("error loading bot", err)
	}
	return &AppointmentBot{bot: b, masterID: masterID}
}

func (a *AppointmentBot) SendNotification(ctx context.Context, message string) {
	if _, err := a.bot.SendMessage(ctx, &tgBot.SendMessageParams{
		ChatID: a.masterID,
		Text:   message,
	}); err != nil {
		log.Printf("error sending message: %v", err)
	}
}

func (a *AppointmentBot) Close(ctx context.Context) (bool, error) {
	return a.bot.Close(ctx)
}
