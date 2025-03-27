package issues

import (
	"gitlab-issues-manager/internal/gitlab-api/constants"
	"io"
	"log"
	"net/http"
)

func CreateIssue(desc, title, gitlabUrl, gitlabToken string) error {
	req, err := http.NewRequest("POST", gitlabUrl, nil)
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
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error while resp.Body.Close()")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	sb := string(body)
	log.Printf("resp body: %v\n", sb)
	return nil
}
