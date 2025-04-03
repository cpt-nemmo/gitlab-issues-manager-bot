package projects

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	constants2 "test/internal/gitlab-api/constants"
	"test/internal/logger"
)

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllProjects(ctx context.Context, gitlabBaseUrl, gitlabToken string) ([]Project, error) {
	l := logger.Enter("gitlab-api.projects.api.GetAllProjects")
	defer func() { logger.Exit(l, "gitlab-api.projects.api.GetAllProjects") }()

	gitLabUrlForGettingAllProjects := gitlabBaseUrl + "/projects"
	req, err := http.NewRequestWithContext(ctx, "GET", gitLabUrlForGettingAllProjects, nil)
	if err != nil {
		log.Printf("[ERROR]: error while POST request: %v\n", err)
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	q := req.URL.Query()
	q.Add("per_page", constants2.PAGINATION)
	req.URL.RawQuery = q.Encode()

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
