package bot

import (
	"context"
	"log"

	tgBot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type NoopBot struct{}

func (n *NoopBot) SendMessage(ctx context.Context, params *tgBot.SendMessageParams) (*models.Message, error) {
	log.Default().Printf("NoopBot: %s", params.Text)
	return nil, nil
}
func (n *NoopBot) Close(ctx context.Context) (bool, error) {
	return true, nil
}
