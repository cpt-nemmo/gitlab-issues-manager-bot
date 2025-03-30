package views

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-issues-manager/internal/bot"
	"gitlab-issues-manager/internal/bot/chat_types"
	"gitlab-issues-manager/internal/bot/message_sender"
	"log"
)

func ViewCmdStart() bot.ViewFunc {
	return func(b *bot.Bot, update tgbotapi.Update, gitlabBaseUrl, gitlabToken string) error {
		chatType := update.Message.Chat.Type

		switch chatType {
		case chat_types.GROUP_CHAT:
			if err := message_sender.Send("Данная команда доступна только в личке с ботом 😔", b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
			}
			return nil
		case chat_types.SUPERGROUP_CHAT:
			if err := message_sender.Send("Данная команда доступна только в личке с ботом 😔", b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
			}
			return nil
		}
		if err := message_sender.Send("Holla!", b.Api, update); err != nil {
			log.Printf("[ERROR] error while sending text msg: %v", err)
		}
		return nil
	}
}
