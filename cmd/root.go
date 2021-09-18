package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/torbensky/pr-watcher/lib"
)

var owner string
var repo string
var prNum int

var rootCmd = &cobra.Command{
	Use:   "pr-watcher",
	Short: "PR Watcher watches a PR and notifies you of changes",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		ghCreds := os.Getenv("GITHUB_ACCESS_TOKEN")
		ctx := context.Background()
		gh := lib.NewGitHub(ctx, ghCreds)

		var last *lib.RepositoryView
		for {
			resp, err := gh.QueryPRStatus(ctx, owner, repo, prNum)
			if err != nil {
				log.Fatal(err)
			}

			// initialize data
			if last == nil {
				last = resp
				continue
			}

			changes := lib.Compare(last, resp)
			for _, c := range changes {
				switch c {
				case lib.NEW_COMMIT:
					lib.Notify("New Commit", "A new commit has started")
				case lib.REVIEW_CHANGE:
					lib.Notify("PR Reviewed", "New review received")
				case lib.CHECK_FAILURE:
					lib.Notify("PR Check", "PR check status updated")
				}
			}
			last = resp

			time.Sleep(time.Second * 10)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&owner, "owner", "o", "", "Repository owner (organization or user)")
	err := rootCmd.MarkPersistentFlagRequired("owner")
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "Repository name")
	err = rootCmd.MarkPersistentFlagRequired("repo")
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().IntVarP(&prNum, "pullRequest", "p", 0, "PR number")
	err = rootCmd.MarkPersistentFlagRequired("pullRequest")
	if err != nil {
		log.Fatal(err)
	}
}
