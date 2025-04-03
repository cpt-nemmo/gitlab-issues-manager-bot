package views

import (
	tgbotapi "github.com/cpt-nemmo/telegram-bot-api"
	"log"
	"test/internal/bot"
	"test/internal/bot/chat_types"
	"test/internal/bot/message_sender"
	"test/internal/logger"
)

func ViewCmdHelp() bot.ViewFunc {
	l := logger.Enter("bot.views.view_cmd_help.ViewCmdHelp")
	defer func() { logger.Exit(l, "bot.views.view_cmd_help.ViewCmdHelp") }()

	return func(b *bot.Bot, update tgbotapi.Update, gitlabBaseUrl, gitlabToken string) error {
		text := "📖 Инструкция для чайников 📖\n\n\nДанный бот помогает хэндлить сообщения с тэгом <b><i>#issue</i></b>.\n" +
			"Чтобы начать им пользваться вы должны сделать следующее:\n\n" +
			"1. Добавить бота в чатик\n\n" +
			"2. Объявить проект в который вы будете постить ваши исуйки. Сделать это можно" +
			"с помощью команды <i>/setproject</i>.\n\n" +
			"3. Дальше просто пишите сообщения в любом удобном формате. Главное чтобы сообщение " +
			"содержало \n\n<b>'Title:'</b> \n\n<b>'Desc:'</b> " +
			"\n\nЗаголовок исуйки и описание соответственно. И чтобы они шли именно в таком порядке." +
			"Сначала 'Title:' --> после 'Desc:' .\n\n" +
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
