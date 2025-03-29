package bot

import (
	"context"
	"log"

	tgBot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AppointmentBot struct {
	bot      *tgBot.Bot
	masterID int64
}

func New(apikey string, chatID int64) *AppointmentBot {
	middleware := func(next tgBot.HandlerFunc) tgBot.HandlerFunc {
		return func(ctx context.Context, b *tgBot.Bot, update *models.Update) {
			if update.Message.From.ID != chatID ||
				update.Message.Chat.Type != "private" ||
				update.Message.Chat.ID != chatID {
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

	b, err := tgBot.New(apikey, opts...)
	if err != nil {
		log.Fatal("error loading bot", err)
	}
	return &AppointmentBot{bot: b, masterID: chatID}
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
