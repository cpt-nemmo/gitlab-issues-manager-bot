package bot

import (
	"context"
	"fmt"
	"gitlab-issues-manager/internal/bot/message_sender"
	"gitlab-issues-manager/internal/gitlab-api/issues"
	"gitlab-issues-manager/internal/utils"
	"log"
	"runtime/debug"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
}

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{api: api}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}

	b.cmdViews[cmd] = view
}

func (b *Bot) Run(ctx context.Context, gitlabUrl, gitlabToken string) error {
	fmt.Println("Bot has been started")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Minute)
			b.handleUpdate(updateCtx, update, gitlabUrl, gitlabToken)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update, gitlabUrl, gitlabToken string) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()
	msgText := update.Message.Text
	if strings.Contains(msgText, "#issue") {
		issue, err := utils.ParseIssue(msgText)
		if err != nil {
			log.Printf("[ERROR] parse issue: %v", err)
			return
		}
		fmt.Println("issue:", issue)
		err = issues.CreateIssue(
			issue.Description,
			issue.Title,
			gitlabUrl,
			gitlabToken,
		)
		if err != nil {
			log.Printf("[ERROR] while creating issue: %v", err)
			return
		}
	}

	if (update.Message == nil || !update.Message.IsCommand()) && update.CallbackQuery == nil {
		return
	}

	var view ViewFunc

	if !update.Message.IsCommand() {
		return
	}

	cmd := update.Message.Command()

	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		if err := message_sender.Send(
			"Сорян, такой команды нету",
			b.api,
			update); err != nil {
			log.Printf("[ERROR] failed to send message in existing command checking: %v", err)
		}
		return
	}

	view = cmdView

	if err := view(ctx, b.api, update); err != nil {
		log.Printf("[ERROR] failed to execute view: %v", err)

		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error")); err != nil {
			log.Printf("[ERROR] failed to send error message: %v", err)
		}
	}
}
