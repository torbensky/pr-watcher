package prwatcher

import (
	"fmt"
	"reflect"

	"github.com/r3labs/diff/v2"
)

type ChangeType uint

const (
	NONE = iota
	REPO_STATUS
	REVIEW_CHANGE
	CHECK_CHANGE
	NEW_COMMIT
)

func Compare(prev, current *RepositoryView) ([]ChangeType, error) {
	var changes []ChangeType

	changelog, err := diff.Diff(prev, current, diff.TagName("json"), diff.DisableStructValues(), diff.SliceOrdering(false))
	if err != nil {
		return []ChangeType{NONE}, err
	}

	fmt.Println("COMPARING")
	// checkChanges := map[string]string{}
	for _, c := range changelog {
		fmt.Println(c.Path, c.From, c.To)

		if len(c.Path) == 3 && c.Path[2] == "state" {
			changes = append(changes, REPO_STATUS)
			continue
		}

		// Change was to commit status
		if len(c.Path) == 10 && c.Path[9] == "state" {
			changes = append(changes, CHECK_CHANGE)
			continue
		}

		// review change
		if len(c.Path) == 6 && c.Path[2] == "reviews" {
			changes = append(changes, REVIEW_CHANGE)
			continue
		}

		// new commit
		if len(c.Path) == 7 && c.Path[6] == "abbreviatedOid" {
			// when we see a new commit, we ignore other changes
			return []ChangeType{NEW_COMMIT}, nil
		}

	}

	if len(changes) > 0 {
		return changes, nil
	}

	return []ChangeType{NONE}, nil
}

type cDiffer struct {
	DiffFunc (func(path []string, a, b reflect.Value, p interface{}) error)
}

func (o *cDiffer) InsertParentDiffer(dfunc func(path []string, a, b reflect.Value, p interface{}) error) {
	o.DiffFunc = dfunc
}

func (o *cDiffer) Match(a, b reflect.Value) bool {
	return diff.AreType(a, b, reflect.TypeOf(CommitStatusContext{}))
}
func (o *cDiffer) Diff(cl *diff.Changelog, path []string, a, b reflect.Value) error {
	if a.String() == b.String() {
		cl.Add(diff.UPDATE, path, a.Interface(), b.Interface())
	}
	return nil
}
