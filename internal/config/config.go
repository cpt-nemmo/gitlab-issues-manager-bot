package config

import (
	"errors"
	"github.com/joho/godotenv"
	"gitlab-issues-manager/internal/gitlab-api/constants"
	"log"
	"os"
)

type GitLabProjectsConfig struct {
	DefaultProjectId int
	Projects         map[string]string
}

type BotConfig struct {
	Token string
}

type GitlabConfig struct {
	Token   string
	BaseUrl string
}

type AppConfig struct {
	Gitlab         GitlabConfig
	Bot            BotConfig
	GitLabProjects GitLabProjectsConfig
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}
}

func (cfg *AppConfig) LoadAppConfig() error {
	token, present := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !present {
		return errors.New("telegram bot token is not set")
	}
	cfg.Bot.Token = token

	url, present := os.LookupEnv("GITLAB_BASE_URL")
	if !present {
		return errors.New("gitlab base url is not set")
	}
	cfg.Gitlab.BaseUrl = url

	gitlabToken, present := os.LookupEnv("GITLAB_TOKEN")
	if !present {
		return errors.New("gitlab token is not set")
	}
	cfg.Gitlab.Token = gitlabToken

	cfg.GitLabProjects.DefaultProjectId = constants.DEFAULT_PROJECT_ID

	return nil
}
