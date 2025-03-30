package projects

import (
	"context"
	"encoding/json"
	"gitlab-issues-manager/internal/logger"
	"io"
	"log"
	"net/http"
)

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllProjects(ctx context.Context, gitlabBaseUrl, gitlabToken string) ([]Project, error) {
	l := logger.Enter("bot.views.view_cmd_getCurrentProject.GetAllProjects")
	defer func() { logger.Exit(l, "bot.views.view_cmd_getCurrentProject.GetAllProjects") }()

	gitLabUrlForGettingAllProjects := gitlabBaseUrl + "/projects"
	req, err := http.NewRequestWithContext(ctx, "GET", gitLabUrlForGettingAllProjects, nil)
	if err != nil {
		log.Printf("[ERROR]: error while POST request: %v\n", err)
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", gitlabToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR]: error while GET request: %v\n", err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error while resp.Body.Close()")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR]: error while io.ReadAll: %v\n", err)
		return nil, err
	}

	var v []Project

	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
