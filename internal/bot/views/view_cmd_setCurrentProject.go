package views

import (
	"context"
	"fmt"
	tgbotapi "github.com/cpt-nemmo/telegram-bot-api"
	"log"
	"test/internal/bot"
	"test/internal/bot/chat_types"
	"test/internal/bot/constants"
	"test/internal/bot/message_sender"
	constants2 "test/internal/gitlab-api/constants"
	"test/internal/gitlab-api/projects"
	"test/internal/logger"
	"test/internal/utils"
	"time"
)

func ViewCmdSetCurrentProject() bot.ViewFunc {
	l := logger.Enter("bot.views.view_cmd_setCurrentProject.ViewCmdSetCurrentProject")
	defer func() { logger.Exit(l, "bot.views.view_cmd_setCurrentProject.ViewCmdSetCurrentProject") }()

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

		var rows [][]tgbotapi.InlineKeyboardButton
		for _, project := range projectss {
			button := tgbotapi.NewInlineKeyboardButtonData(project.Name, project.Name)
			row := tgbotapi.NewInlineKeyboardRow(button)
			rows = append(rows, row)
		}
		firstMenuMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

		var user = update.Message.From.UserName
		var textForKeyboard = fmt.Sprintf("<b>@%v</b>, выбери проект ниже", user)
		messageID, err := message_sender.SendMenu(b.Api, firstMenuMarkup, textForKeyboard, update.Message.Chat.ID, update.Message.MessageThreadID)
		if err != nil {
			return err
		}
		b.ChatState = constants.PENDING_KEYBOARD_ANSWER
		b.MessageIdForDeletion = messageID
		b.ChatID = update.Message.Chat.ID
		return nil
	}
}
