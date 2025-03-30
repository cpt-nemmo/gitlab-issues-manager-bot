package views

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-issues-manager/internal/bot"
	"gitlab-issues-manager/internal/bot/chat_types"
	"gitlab-issues-manager/internal/bot/constants"
	"gitlab-issues-manager/internal/bot/message_sender"
	constants2 "gitlab-issues-manager/internal/gitlab-api/constants"
	"gitlab-issues-manager/internal/gitlab-api/projects"
	"gitlab-issues-manager/internal/utils"
	"log"
	"time"
)

func ViewCmdSetCurrentProject() bot.ViewFunc {
	return func(b *bot.Bot, update tgbotapi.Update, gitlabBaseUrl, gitlabToken string) error {

		chatType := update.Message.Chat.Type

		switch chatType {
		case chat_types.PRIVATE_CHAT:
			if err := message_sender.Send("Данная команда доступна только в беседе 😔", b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
			}
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), constants2.TIMEOUT*time.Second)
		defer cancel()
		if err := message_sender.Send("⏳ Делаем запрос на гитлаб...", b.Api, update); err != nil {
			return err
		}
		projectss, err := projects.GetAllProjects(ctx, gitlabBaseUrl, gitlabToken)
		if err != nil {
			err = message_sender.Send("Сетевая ошибка при запросе на ваш гитлаб. Соре... 😔", b.Api, update)
			return err
		}
		b.Projects = utils.ConvertFromSliceToMap(projectss)

		var rows [][]tgbotapi.KeyboardButton
		for _, project := range projectss {
			button := tgbotapi.NewKeyboardButton(project.Name)
			row := tgbotapi.NewKeyboardButtonRow(button)
			rows = append(rows, row)
		}
		firstMenuMarkup := tgbotapi.NewReplyKeyboard(rows...)

		var user = update.Message.From.UserName
		var textForKeyboard = fmt.Sprintf("<b>@%v</b>, выбери проект на клавиатуре:", user)
		if err := message_sender.SendMenu(b.Api, firstMenuMarkup, textForKeyboard, update.Message.Chat.ID); err != nil {
			return err
		}
		b.ChatState = constants.PENDING_KEYBOARD_ANSWER
		return nil
	}
}
