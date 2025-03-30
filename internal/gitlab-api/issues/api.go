package issues

import (
	"context"
	"fmt"
	"gitlab-issues-manager/internal/gitlab-api/constants"
	"gitlab-issues-manager/internal/logger"
	"io"
	"log"
	"net/http"
)

func CreateIssue(ctx context.Context, desc, title, gitlabBaseUrl, gitlabToken string, currentProjID int) error {
	l := logger.Enter("gitlab-api.issues.api.CreateIssue")
	defer func() { logger.Exit(l, "gitlab-api.issues.api.CreateIssue") }()

	gitLabUrlForAddingIssues := gitlabBaseUrl + fmt.Sprintf("/projects/%v/issues", currentProjID)
	req, err := http.NewRequestWithContext(ctx, "POST", gitLabUrlForAddingIssues, nil)
	if err != nil {
		log.Printf("[ERROR]: error while POST request: %v\n", err)
		return err
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
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error while resp.Body.Close()")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR]: error while io.ReadAll: %v\n", err)
		return err
	}
	log.Printf("[RESPONSE BODY FROM GITLAB]: %v\n", string(body))
	return nil
}
