package issues

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"test/internal/gitlab-api/constants"
	"test/internal/logger"
)

type GitlabResponse struct {
	Url string `json:"web_url"`
}

type Counts struct {
	All    float64 `json:"all"`
	Closed float64 `json:"closed"`
	Opened float64 `json:"opened"`
}

type Statistics struct {
	Counts Counts `json:"counts"`
}

type IssueStatistics struct {
	Statistics Statistics `json:"statistics"`
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

func GetStatisticByProjectID(ctx context.Context, gitlabBaseUrl, gitlabToken string, currentProjID int) (float64, float64, float64, error) {
	l := logger.Enter("gitlab-api.issues.api.GetStatisticByProjectID")
	defer func() { logger.Exit(l, "gitlab-api.issues.api.GetStatisticByProjectID") }()

	gitLabUrlForAddingIssues := gitlabBaseUrl + fmt.Sprintf("/projects/%v/issues_statistics", currentProjID)
	req, err := http.NewRequestWithContext(ctx, "GET", gitLabUrlForAddingIssues, nil)
	if err != nil {
		log.Printf("[ERROR]: error while POST request: %v\n", err)
		return 0, 0, 0, err
	}
	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR]: error while POST request: %v\n", err)
		return 0, 0, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error while resp.Body.Close()")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR]: error while io.ReadAll: %v\n", err)
		return 0, 0, 0, err
	}
	log.Printf("[RESPONSE BODY FROM GITLAB]: %v\n", string(body))

	var v IssueStatistics

	err = json.Unmarshal(body, &v)
	if err != nil {
		return 0, 0, 0, err
	}

	//fmt.Println("ISSUE STATISTICS: ",
	//	v.Statistics.Counts.All,
	//	v.Statistics.Counts.Opened,
	//	v.Statistics.Counts.Closed,
	//)

	return v.Statistics.Counts.All, v.Statistics.Counts.Opened, v.Statistics.Counts.Closed, nil
}
