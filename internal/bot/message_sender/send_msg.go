package message_sender

import (
	tgbotapi "github.com/cpt-nemmo/telegram-bot-api"
	"test/internal/bot/markup_formatter"
)

func SendMenu(bot *tgbotapi.BotAPI, firstMenuMarkup tgbotapi.InlineKeyboardMarkup, textForKeyboard string, chatId int64, threadID int) (int, error) {
	msg := tgbotapi.NewMessage(chatId, threadID, textForKeyboard)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	sentMsg, err := bot.Send(msg)
	if err != nil {
		return 0, err
	}
	messageID := sentMsg.MessageID

	return messageID, nil
}

func Send(text string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.MessageThreadID, markup_formatter.Replacer(text))
	reply.ParseMode = "MarkdownV2"
	_, err := bot.Send(reply)
	return err
}

func SendHTML(text string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.MessageThreadID, text)
	reply.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(reply)
	return err
}
