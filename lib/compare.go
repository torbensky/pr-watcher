package lib

import "fmt"

// Describes a change to the PR
type Change struct {
	// A machine-readable change-type
	Type ChangeType
	// A human-readable summary of the change
	Summary string
}

func (c Change) String() string {
	return c.Summary
}

type ChangeType uint

const (
	REVIEW_CHANGE = iota
	CHECK_FAILURE
	NEW_COMMIT
	ALL_CHECKS_SUCCESS
)

func Compare(prev, current *RepositoryView) []Change {
	var changes []Change

	prevPR := prev.Repository.PullRequest
	currentPR := current.Repository.PullRequest
	prevCommit := prevPR.Commits.Nodes[0]
	currentCommit := currentPR.Commits.Nodes[0]

	// check if the commit changed
	// if it did, that's the only change we should return since it changes everything
	if prevCommit.Commit.AbbreviatedOID != currentCommit.Commit.AbbreviatedOID {
		return []Change{{NEW_COMMIT, "PR updated with a new commit 🏗️"}}
	}

	// check if a review changed
	pastReviewStates := map[string]string{}
	for _, r := range prevPR.Reviews.Nodes {
		pastReviewStates[r.Author.Login] = r.State
	}
	for _, r := range currentPR.Reviews.Nodes {
		if _, ok := pastReviewStates[r.Author.Login]; !ok {
			// this review is new
			changes = append(changes, Change{REVIEW_CHANGE, fmt.Sprintf("%s %s your PR", r.Author.Login, r.State)})
		}
	}

	// check if a PR status check changed
	pastStatusChecks := map[string]StatusState{}
	for _, c := range prevCommit.Commit.Status.Contexts {
		pastStatusChecks[c.ID] = c.State
	}
	unsuccessfullChecks := len(currentCommit.Commit.Status.Contexts)
	checkChange := false
	for _, c := range currentCommit.Commit.Status.Contexts {

		// determine whether this check has changed
		newStatus := false
		if lastStatus, ok := pastStatusChecks[c.ID]; ok {
			if lastStatus != c.State {
				// this check had its state change, so that's new
				newStatus = true
			}
		} else {
			// this check was not seen before, so it's new
			newStatus = true
		}

		// later we need to know if any checks changed
		checkChange = checkChange || newStatus

		switch c.State {
		case Success:
			unsuccessfullChecks--
		case Failure:
			fallthrough
		case Error:
			// failures or errors are treated the same
			if newStatus {
				changes = append(changes, Change{CHECK_FAILURE, fmt.Sprintf("%s %s 🔥", c.Context, c.State)})
			}
		}
	}

	// if all checks have passed
	if checkChange && unsuccessfullChecks == 0 {
		changes = append(changes, Change{ALL_CHECKS_SUCCESS, "All status checks successfull. Your PR is ready to go 🚀"})
	}

	return changes
}
