package message_sender

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-issues-manager/internal/bot/markup_formatter"
)

func Send(text string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	var (
		msgText = fmt.Sprintf(
			text,
		)
		reply = tgbotapi.NewMessage(update.Message.Chat.ID, markup_formatter.Replacer(msgText))
	)

	reply.ParseMode = "MarkdownV2"

	if _, err := bot.Send(reply); err != nil {
		return err
	}
	return nil
}
