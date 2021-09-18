package lib

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
	"golang.org/x/oauth2"
)

type GitHub struct {
	client *graphql.Client
}

func NewGitHub(ctx context.Context, accessToken string) *GitHub {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	oauthClient := oauth2.NewClient(ctx, ts)
	oauthClient.Timeout = time.Second * 10
	gqlClient := graphql.NewClient("https://api.github.com/graphql", graphql.WithHTTPClient(oauthClient))

	return &GitHub{
		client: gqlClient,
	}
}

func (gh GitHub) QueryPRStatus(ctx context.Context, owner, repo string, pr int) (*RepositoryView, error) {
	req := graphql.NewRequest(prStateQuery)
	req.Var("owner", owner)
	req.Var("repo", repo)
	req.Var("pull_number", pr)

	var result RepositoryView
	err := gh.client.Run(ctx, req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

const prStateQuery = `query($owner: String!, $repo: String!, $pull_number: Int!) {
	repository(owner: $owner, name:$repo) {
	  pullRequest(number:$pull_number) {
		title
		state
		mergeable
		reviews(last: 1) {
			nodes {
			state
			author {
				login
			}
			}
		}
		commits(last: 1) {
			nodes {
				commit {
					abbreviatedOid
					status {
						contexts {
							id
							context
							description
							state
						}
					}
				}
			}
		}
	}
}
}`

type RepositoryView struct {
	Repository struct {
		PullRequest struct {
			Title     string `json:"title"`
			State     string `json:"state"`
			Mergeable string `json:"mergeable"`
			Reviews   *struct {
				Nodes []ReviewNode `json:"nodes"`
			} `json:"reviews"`
			Commits struct {
				Nodes []struct {
					Commit struct {
						AbbreviatedOID string `json:"abbreviatedOid"`
						Status         struct {
							Contexts []CommitStatusContext `json:"contexts"`
						} `json:"status"`
					} `json:"commit"`
				} `json:"nodes"`
			} `json:"commits"`
		} `json:"pullRequest"`
	} `json:"repository"`
}

type ReviewNode struct {
	State  string `json:"state"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}

type StatusState string

const (
	Expected StatusState = "EXPECTED"
	Error    StatusState = "ERROR"
	Failure  StatusState = "FAILURE"
	Pending  StatusState = "PENDING"
	Success  StatusState = "SUCCESS"
)

type CommitStatusContext struct {
	Context     string      `json:"context"`
	Description string      `json:"description"`
	State       StatusState `json:"state"`
	ID          string      `json:"id"`
}
