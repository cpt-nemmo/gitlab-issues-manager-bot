package views

import (
	tgbotapi "github.com/cpt-nemmo/telegram-bot-api"
	"log"
	"test/internal/bot"
	"test/internal/bot/chat_types"
	"test/internal/bot/message_sender"
	"test/internal/logger"
)

func ViewCmdStart() bot.ViewFunc {
	l := logger.Enter("bot.views.view_cmd_start.ViewCmdStart")
	defer func() { logger.Exit(l, "bot.views.view_cmd_start.ViewCmdStart") }()

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
