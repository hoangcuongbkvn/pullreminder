package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

type params map[string]interface{}

// GetPRs get list pullrequest
// set multiple repositories by comma
// ex: pkg-booking-core,pkg-booking-restaurant
func GetPRs(token, owner, repo string) []*github.PullRequest {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var results []*github.PullRequest
	repos := strings.Split(repo, ",")

	for _, r := range repos {
		result, _, _ := client.PullRequests.List(ctx, owner, r, &github.PullRequestListOptions{})
		results = append(results, result...)
	}

	return results
}

// SendToSlack notify to Slack
func SendToSlack(WebhoolURL string, text io.Reader) {
	client := http.DefaultClient
	client.Post(WebhoolURL, "application/json", text)
}

func main() {

	prs := GetPRs(os.Getenv("GITHUB_API_KEY"), os.Getenv("GITHUB_OWNER"), os.Getenv("GITHUB_REPO"))

	var data []params
	for _, pr := range prs {
		data = append(data, params{"fallback": *pr.Title,
			"pretext":    *pr.Title,
			"title":      *pr.Title,
			"title_link": *pr.HTMLURL,
			"text":       *pr.Body,
			"color":      "#ff9000"})

	}

	printData, _ := json.Marshal(params{"attachments": data})

	SendToSlack(os.Getenv("SLACK_WEBHOOK"), bytes.NewBuffer(printData))

}
