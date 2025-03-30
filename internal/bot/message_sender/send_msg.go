package message_sender

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-issues-manager/internal/bot/markup_formatter"
)

func DeleteMenu(bot *tgbotapi.BotAPI, textForKeyboard string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, textForKeyboard)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(msg)
	return err
}

func SendMenu(bot *tgbotapi.BotAPI, firstMenuMarkup tgbotapi.ReplyKeyboardMarkup, textForKeyboard string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, textForKeyboard)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

func Send(text string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, markup_formatter.Replacer(text))
	reply.ParseMode = "MarkdownV2"
	_, err := bot.Send(reply)
	return err
}

func SendHTML(text string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	reply.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(reply)
	return err
}
