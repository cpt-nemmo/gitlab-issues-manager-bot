package views

import (
	"fmt"
	tgbotapi "github.com/cpt-nemmo/telegram-bot-api"
	"log"
	"test/internal/bot"
	"test/internal/bot/chat_types"
	"test/internal/bot/message_sender"
	"test/internal/logger"
)

func ViewCmdGetCurrentProject() bot.ViewFunc {
	l := logger.Enter("bot.views.view_cmd_getCurrentProject.ViewCmdGetCurrentProject")
	defer func() { logger.Exit(l, "bot.views.view_cmd_getCurrentProject.ViewCmdGetCurrentProject") }()

	return func(b *bot.Bot, update tgbotapi.Update, gitlabBaseUrl, gitlabToken string) error {
		chatType := update.Message.Chat.Type

		switch chatType {
		case chat_types.PRIVATE_CHAT:
			if err := message_sender.Send("Данная команда доступна только в беседе 😔", b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
			}
			return nil
		}

		currentProject := b.CurrentProject.Name
		if currentProject == "" {
			if err := message_sender.Send("⛔️Требуется сначала установить проект!⛔️", b.Api, update); err != nil {
				return err
			}
			return nil
		}

		msgText := fmt.Sprintf("🌚Текущий проект:\n <b>%s</b>", currentProject)

		if err := message_sender.SendHTML(msgText, b.Api, update); err != nil {
			log.Printf("[ERROR] error while sending text msg: %v", err)
		}
		return nil
	}
}
