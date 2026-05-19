package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type GithubEvent struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Actor struct {
		ID           int    `json:"id"`
		Login        string `json:"login"`
		DisplayLogin string `json:"display_login"`
		GravatarID   string `json:"gravatar_id"`
		URL          string `json:"url"`
		AvatarURL    string `json:"avatar_url"`
	} `json:"actor"`
	Repo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repo"`
	Payload struct {
		Action string `json:"action"`
	} `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
	Org       struct {
		ID         int    `json:"id"`
		Login      string `json:"login"`
		GravatarID string `json:"gravatar_id"`
		URL        string `json:"url"`
		AvatarURL  string `json:"avatar_url"`
	} `json:"org,omitempty"`
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

	var events []GithubEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Printf("Error parsing API response JSON: %v\n", err)
		return
	}

	if len(events) == 0 {
		fmt.Printf("No recent public activity found for user: %s\n", username)
		return
	}

	fmt.Printf("Events: %+v", events)

}
