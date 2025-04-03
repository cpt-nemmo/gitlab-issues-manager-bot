package views

import (
	"bytes"
	"context"
	"fmt"
	tgbotapi "github.com/cpt-nemmo/telegram-bot-api"
	"github.com/wcharczuk/go-chart"
	"log"
	"test/internal/bot"
	"test/internal/bot/chat_types"
	"test/internal/bot/message_sender"
	constants2 "test/internal/gitlab-api/constants"
	"test/internal/gitlab-api/issues"
	"test/internal/logger"
	"time"
)

func ViewCmdStatistics() bot.ViewFunc {
	l := logger.Enter("bot.views.view_cmd_statistics.ViewCmdStatistics")
	defer func() { logger.Exit(l, "bot.views.view_cmd_statistics.ViewCmdStatistics") }()

	return func(b *bot.Bot, update tgbotapi.Update, gitlabBaseUrl, gitlabToken string) error {

		chatType := update.Message.Chat.Type

		switch chatType {
		case chat_types.PRIVATE_CHAT:
			if err := message_sender.Send("Данная команда доступна только в беседе 😔", b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
			}
			return nil
		}

		if b.CurrentProject.Id == 0 {
			if err := message_sender.Send("⛔️Требуется сначала установить проект!⛔️", b.Api, update); err != nil {
				return err
			}
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), constants2.TIMEOUT*time.Second)
		defer cancel()
		if err := message_sender.Send("⏳ Делаем запрос на гитлаб...", b.Api, update); err != nil {
			return err
		}
		allIssues, openedIssues, closedIssues, err := issues.GetStatisticByProjectID(
			ctx,
			gitlabBaseUrl,
			gitlabToken,
			b.CurrentProject.Id,
		)
		if err != nil {
			log.Printf("[ERROR] error while getting statistics: %v", err)
			return err
		}

		if openedIssues == 0 || closedIssues == 0 {
			if err := message_sender.Send("Проект только начался или ведется очень плохо. Релевантной статистики нету...😔", b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
				return err
			}
			return nil
		}

		openedIssuesProcent := openedIssues / allIssues * 100
		closedIssuesProcent := closedIssues / allIssues * 100

		err = message_sender.SendHTML(
			fmt.Sprintf("📊Статистика по проекту:\n <b>%v</b>\n\n 🔊<i>Всего задач:</i> %v\n\n 🔊<i>Из них готово:</i> %v\n\n", b.CurrentProject.Name, allIssues, closedIssues),
			b.Api, update)

		pie := chart.PieChart{
			Width:  512,
			Height: 512,
			Values: []chart.Value{
				{Value: float64(openedIssues), Label: fmt.Sprintf("Active :: %.2f%", openedIssuesProcent)},
				{Value: float64(closedIssues), Label: fmt.Sprintf("Closed :: %.2f%", closedIssuesProcent)},
			},
		}

		// Рендер в буфер
		buf := bytes.NewBuffer([]byte{})
		if err := pie.Render(chart.PNG, buf); err != nil {
			return err
		}

		// Отправка пользователю
		msg := tgbotapi.NewPhoto(update.Message.Chat.ID, update.Message.MessageThreadID, tgbotapi.FileBytes{
			Name:  "pie.png",
			Bytes: buf.Bytes(),
		})
		if _, err := b.Api.Send(msg); err != nil {
			return err
		}

		if err != nil {
			log.Printf("[ERROR] while creating issue: %v", err)
			err = message_sender.Send("Сетевая ошибка при запросе на ваш гитлаб. Соре... 😔", b.Api, update)
			return err
		}
		return nil
	}
}
