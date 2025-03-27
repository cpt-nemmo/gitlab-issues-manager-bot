package views

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-issues-manager/internal/bot"
	"gitlab-issues-manager/internal/bot/message_sender"
	"log"
)

func ViewCmdStart() bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if err := message_sender.Send("Holla!", bot, update); err != nil {
			log.Printf("[ERROR] error while sending text msg: %v", err)
		}
		return nil
	}
}
