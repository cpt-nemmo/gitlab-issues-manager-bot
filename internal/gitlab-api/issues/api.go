package issues

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab-issues-manager/internal/gitlab-api/constants"
	"gitlab-issues-manager/internal/logger"
	"io"
	"log"
	"net/http"
)

type GitlabResponse struct {
	Url string `json:"web_url"`
}

func CreateIssue(ctx context.Context, desc, title, gitlabBaseUrl, gitlabToken string, currentProjID int) (string, error) {
	l := logger.Enter("gitlab-api.issues.api.CreateIssue")
	defer func() { logger.Exit(l, "gitlab-api.issues.api.CreateIssue") }()

	gitLabUrlForAddingIssues := gitlabBaseUrl + fmt.Sprintf("/projects/%v/issues", currentProjID)
	req, err := http.NewRequestWithContext(ctx, "POST", gitLabUrlForAddingIssues, nil)
	if err != nil {
		log.Printf("[ERROR]: error while POST request: %v\n", err)
		return "", err
	}
	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	q := req.URL.Query()
	q.Add("labels", constants.LABELS)
	q.Add("description", desc)
	q.Add("title", title)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR]: error while POST request: %v\n", err)
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error while resp.Body.Close()")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR]: error while io.ReadAll: %v\n", err)
		return "", err
	}
	log.Printf("[RESPONSE BODY FROM GITLAB]: %v\n", string(body))

	var v GitlabResponse

	err = json.Unmarshal(body, &v)
	if err != nil {
		return "", err
	}

	return v.Url, nil
}
