package views

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-issues-manager/internal/bot"
	"gitlab-issues-manager/internal/bot/chat_types"
	"gitlab-issues-manager/internal/bot/message_sender"
	"log"
)

func ViewCmdGetCurrentProject() bot.ViewFunc {
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
