package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type GitHubEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Repo      struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Action      string `json:"action,omitempty"`   // e.g., "opened", "merged", "created"
	Ref         string `json:"ref,omitempty"`      // e.g., "refs/heads/fix-cron"
	RefType     string `json:"ref_type,omitempty"` // e.g., "branch"
	Number      int    `json:"number,omitempty"`   // Issue or PR number
	PullRequest *struct {
		HTMLURL string `json:"html_url"`
	} `json:"pull_request,omitempty"`
	Issue *struct {
		Title string `json:"title"`
	} `json:"issue,omitempty"`
	Comment *struct {
		Body string `json:"body"`
	} `json:"comment,omitempty"`
}

func main() {
	args := os.Args

	if len(args) > 2 {
		log.Fatal("Usage: github-activity <username>")
	}

	username := args[1]
	fetchGithubActivity(username)
}

func fetchGithubActivity(username string) {
	githubLink := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	req, err := http.NewRequest("GET", githubLink, nil)
	if err != nil {
		fmt.Printf("Error creating new request: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", "go-github-activity-cli")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Cannot retrieve data from github: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var events []GitHubEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Printf("Error parsing API response JSON: %v\n", err)
		return
	}

	if len(events) == 0 {
		fmt.Printf("No recent public activity found for user: %s\n", username)
		return
	}

	printActivityLog(events)
}

func printActivityLog(events []GitHubEvent) {
	for _, event := range events {
		dateStr := event.CreatedAt.Format("2006-01-02 15:04")
		repoName := event.Repo.Name

		switch event.Type {
		case "PushEvent":
			fmt.Printf("[%s] Pushed updates to branch %s at %s\n", dateStr, event.Payload.Ref, repoName)

		case "PullRequestEvent":
			action := event.Payload.Action // "opened" or "merged"
			prNum := event.Payload.Number
			fmt.Printf("[%s] Pull Request #%d %s in %s\n", dateStr, prNum, action, repoName)

		case "IssueCommentEvent":
			fmt.Printf("[%s] Commented on an issue/PR in %s\n", dateStr, repoName)
			if event.Payload.Comment != nil {
				fmt.Printf("  - Context: %s\n", event.Payload.Comment.Body)
			}

		case "CreateEvent":
			fmt.Printf("[%s] Created %s '%s' in %s\n", dateStr, event.Payload.RefType, event.Payload.Ref, repoName)

		case "WatchEvent":
			fmt.Printf("[%s] Starred %s\n", dateStr, repoName)

		case "ForkEvent":
			fmt.Printf("[%s] Forked %s\n", dateStr, repoName)

		default:
			fmt.Printf("[%s] %s occurred in %s\n", dateStr, event.Type, repoName)
		}
	}
}
