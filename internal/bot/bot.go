package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/cpt-nemmo/telegram-bot-api"
	"log"
	"strings"
	"test/internal/bot/chat_types"
	"test/internal/bot/constants"
	"test/internal/bot/message_sender"
	constants2 "test/internal/gitlab-api/constants"
	"test/internal/gitlab-api/issues"
	"test/internal/logger"
	"test/internal/utils"
	"time"
)

type ViewFunc func(b *Bot, update tgbotapi.Update, gitlabBaseUrl, gitlabToken string) error

type CurrProj struct {
	Name string
	Id   int
}

type Bot struct {
	Api                  *tgbotapi.BotAPI
	CmdViews             map[string]ViewFunc
	Projects             map[string]int
	ChatState            string
	CurrentProject       CurrProj
	MessageIdForDeletion int
	ChatID               int64
}

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{Api: api}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.CmdViews == nil {
		b.CmdViews = make(map[string]ViewFunc)
	}

	b.CmdViews[cmd] = view
}

func (b *Bot) Run(
	ctx context.Context,
	gitlabBaseUrl,
	gitlabToken string,
) error {
	fmt.Println("Bot has been started")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.Api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			b.handleUpdate(update, gitlabBaseUrl, gitlabToken)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func handleText(
	b *Bot,
	update tgbotapi.Update,
	gitlabBaseUrl,
	gitlabToken string,
) error {
	l := logger.Enter("bot.bot.handleText")
	defer func() { logger.Exit(l, "bot.bot.handleText") }()

	text := update.Message.Text
	chatType := update.Message.Chat.Type

	switch chatType {
	case chat_types.PRIVATE_CHAT:
		if err := message_sender.Send("В личке с ботом ты не можешь ничего обсуждать. Соре.. 💃", b.Api, update); err != nil {
			log.Printf("[ERROR] error while sending text msg: %v", err)
		}
		return nil
	}

	switch b.ChatState {
	case constants.DEFAULT_CHAT_STATE:
		if strings.Contains(text, "#issue") {
			issue, err := utils.ParseIssue(text)
			if err != nil {
				log.Printf("[ERROR] parse issue: %v", err)
				err = message_sender.SendHTML("❗<b>ОШИБКА</b>❗\n️Проверьте - совпадает ли ваше сообшение правилам оформления исуйки в тг.", b.Api, update)
				return err
			}

			if b.CurrentProject.Id == 0 {
				if err := message_sender.Send("⛔️Требуется сначала установить проект!⛔️", b.Api, update); err != nil {
					return err
				}
				return nil
			}

			ctx, cancel := context.WithTimeout(context.Background(), constants2.TIMEOUT*time.Second)
			defer cancel()
			err = message_sender.Send("⏳ Делаем запрос на гитлаб...", b.Api, update)
			url, err := issues.CreateIssue(
				ctx,
				issue.Description,
				issue.Title,
				gitlabBaseUrl,
				gitlabToken,
				b.CurrentProject.Id,
			)
			if err != nil {
				log.Printf("[ERROR] while creating issue: %v", err)
				err = message_sender.Send("Сетевая ошибка при запросе на ваш гитлаб. Соре... 😔", b.Api, update)
				return err
			}
			if err := message_sender.SendHTML(fmt.Sprintf("☑️ Исуйка была записана в проект: <b><i>%v</i></b> \n\nСсылка на нее: %v", b.CurrentProject.Name, url), b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
			}
		}
	}

	return nil
}

func handleCommand(b *Bot, update tgbotapi.Update, gitlabBaseUrl, gitlabToken string) error {
	l := logger.Enter("bot.bot.handleCommand")
	defer func() { logger.Exit(l, "bot.bot.handleCommand") }()

	var view ViewFunc
	var command string

	chatType := update.Message.Chat.Type
	switch chatType {
	case chat_types.GROUP_CHAT:
		if !strings.Contains(update.Message.Text, "@") {
			if err := message_sender.Send("Такой команды не существует 😔", b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
				return err
			}
			return nil
		}
		rowCommand := update.Message.Text
		startInx := strings.Index(rowCommand, "@")
		command = strings.TrimSpace(rowCommand[:startInx])
	case chat_types.SUPERGROUP_CHAT:
		if !strings.Contains(update.Message.Text, "@") {
			if err := message_sender.Send("Такой команды не существует 😔", b.Api, update); err != nil {
				log.Printf("[ERROR] error while sending text msg: %v", err)
				return err
			}
			return nil
		}
		rowCommand := update.Message.Text
		startInx := strings.Index(rowCommand, "@")
		command = strings.TrimSpace(rowCommand[:startInx])
	case chat_types.PRIVATE_CHAT:
		command = update.Message.Text
	}

	cmdView, ok := b.CmdViews[command]
	if !ok {
		if err := message_sender.Send(
			"Сорян, такой команды нету",
			b.Api,
			update); err != nil {
			log.Printf("[ERROR] failed to send message in existing command checking: %v", err)
		}
		return nil
	}

	view = cmdView

	if err := view(b, update, gitlabBaseUrl, gitlabToken); err != nil {
		log.Printf("[ERROR] failed to execute view: %v", err)

		if _, err := b.Api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.MessageThreadID, "Internal error")); err != nil {
			log.Printf("[ERROR] failed to send error message: %v", err)
		}
	}
	return nil
}

func handleCallback(b *Bot, update tgbotapi.Update) {
	l := logger.Enter("bot.bot.handleCallback")
	defer func() { logger.Exit(l, "bot.bot.handleCallback") }()

	switch b.ChatState {
	case constants.DEFAULT_CHAT_STATE:
		break
	case constants.PENDING_KEYBOARD_ANSWER:
		callbackData := update.CallbackQuery.Data
		chatID := b.ChatID
		messageID := b.MessageIdForDeletion
		_, exist := b.Projects[callbackData]
		if exist {
			b.CurrentProject.Id = b.Projects[callbackData]
			b.CurrentProject.Name = callbackData

			msg := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0)})
			_, err := b.Api.Send(msg)
			if err != nil {
				return
			}
			b.ChatState = constants.DEFAULT_CHAT_STATE
			edit := tgbotapi.NewEditMessageText(
				chatID,
				messageID,
				fmt.Sprintf("🌟Вы выбрали проект:\n <b>%v</b>", b.CurrentProject.Name),
			)
			edit.ParseMode = tgbotapi.ModeHTML
			_, err = b.Api.Send(edit)
			if err != nil {
				return
			}
		}
	}
}

func handleMessage(
	b *Bot,
	update tgbotapi.Update,
	gitlabBaseUrl, gitlabToken string,
) {
	l := logger.Enter("bot.bot.handleMessage")
	defer func() { logger.Exit(l, "bot.bot.handleMessage") }()

	user := update.Message.From
	text := update.Message.Text

	if user == nil {
		return
	}

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(b, update, gitlabBaseUrl, gitlabToken)
	} else {
		err = handleText(b, update, gitlabBaseUrl, gitlabToken)
	}

	if err != nil {
		log.Printf("An error occured while handling message: %s", err.Error())
	}
}

func (b *Bot) handleUpdate(
	update tgbotapi.Update,
	gitlabBaseUrl, gitlabToken string,
) {
	l := logger.Enter("bot.bot.handleUpdate")
	defer func() { logger.Exit(l, "bot.bot.handleUpdate") }()

	switch {
	case update.Message != nil:
		if update.Message.Chat.ID < 0 && update.Message.Chat.ID != -1002576766431 {
			if err := message_sender.Send("⛔️Данный чат не авторизован!⛔️", b.Api, update); err != nil {
				return
			}
		} else {
			handleMessage(b, update, gitlabBaseUrl, gitlabToken)
			break
		}
	case update.CallbackQuery != nil:
		handleCallback(b, update)
	}
}
