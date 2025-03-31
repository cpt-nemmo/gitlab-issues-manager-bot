package di

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-issues-manager/internal/bot"
	"gitlab-issues-manager/internal/bot/views"
	"gitlab-issues-manager/internal/config"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type DI struct {
	appCfg *config.AppConfig
	Logger *zap.Logger
}

func (di *DI) Init() error {
	err := di.loadConfig()
	if err != nil {
		return err
	}

	return nil
}

func (di *DI) loadConfig() error {
	di.appCfg = &config.AppConfig{}
	err := di.appCfg.LoadAppConfig()
	if err != nil {
		return errors.New("failed to load app config: " + err.Error())
	}

	return nil
}

func (di *DI) StartBot() {
	botApi, err := tgbotapi.NewBotAPI(di.appCfg.Bot.Token)

	if err != nil {
		log.Printf("failed to create bot: %v", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	gitlabIssuesManagerBot := bot.New(botApi)
	gitlabIssuesManagerBot.RegisterCmdView(
		"/start",
		views.ViewCmdStart(),
	)
	gitlabIssuesManagerBot.RegisterCmdView(
		"/setproject",
		views.ViewCmdSetCurrentProject(),
	)
	gitlabIssuesManagerBot.RegisterCmdView(
		"/getproject",
		views.ViewCmdGetCurrentProject(),
	)
	gitlabIssuesManagerBot.RegisterCmdView(
		"/help",
		views.ViewCmdHelp(),
	)
	gitlabIssuesManagerBot.RegisterCmdView(
		"/statistic",
		views.ViewCmdStatistics(),
	)

	if err := gitlabIssuesManagerBot.Run(
		ctx, di.appCfg.Gitlab.BaseUrl,
		di.appCfg.Gitlab.Token,
	); err != nil {
		log.Printf("[ERROR] failed to run bot: %v", err)
	}
}
