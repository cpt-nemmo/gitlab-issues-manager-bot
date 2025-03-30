package views

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-issues-manager/internal/bot"
	"gitlab-issues-manager/internal/bot/chat_types"
	"gitlab-issues-manager/internal/bot/message_sender"
	"log"
)

func ViewCmdHelp() bot.ViewFunc {
	return func(b *bot.Bot, update tgbotapi.Update, gitlabBaseUrl, gitlabToken string) error {
		text := "Данный бот помогает хэндлить сообщения с тэгом <b><i>#issue</i></b>.\n" +
			"Чтобы начать им пользваться вы должны сделать следующее:\n\n" +
			"1. Добавить бота в чатик\n" +
			"2. Объявить проект в который вы будете постить ваши исуйки. Сделать это можно" +
			"с помощью команды <i>/setproject</i>.\n" +
			"3. Дальше просто пишите сообщения в любом удобном формате. Главное чтобы сообщение " +
			"содержало \n\n'<b><i>title:</i></b>' \n\n'<b><i>desc:</i></b>'. \nЗаголовок исуйки и описание соответственно.\n" +
			"4. Ну и глянуть какой сейчас установлен проект можно командой: <i>/getproject</i>."

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
		if err := message_sender.SendHTML(text, b.Api, update); err != nil {
			log.Printf("[ERROR] error while sending text msg: %v", err)
		}
		return nil
	}
}
