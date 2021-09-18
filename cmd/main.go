package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	prwatcher "github.com/torbensky/pr-watcher"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ghCreds := os.Getenv("GITHUB_ACCESS_TOKEN")
	ctx := context.Background()
	gh := prwatcher.NewGitHub(ctx, ghCreds)

	var last *prwatcher.RepositoryView
	for {
		resp, err := gh.QueryPRStatus(ctx, "techdroplabs", "dyspatch", 6554)
		if err != nil {
			log.Fatal(err)
		}

		// initialize data
		if last == nil {
			last = resp
			continue
		}

		changes := prwatcher.Compare(last, resp)
		for _, c := range changes {
			switch c {
			case prwatcher.NEW_COMMIT:
				prwatcher.Notify("New Commit", "A new commit has started")
			case prwatcher.REVIEW_CHANGE:
				prwatcher.Notify("PR Reviewed", "New review received")
			case prwatcher.CHECK_FAILURE:
				prwatcher.Notify("PR Check", "PR check status updated")
			}
		}
		last = resp

		time.Sleep(time.Second * 10)
	}
}
