package lib_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torbensky/pr-watcher/lib"
)

func TestCompare(t *testing.T) {
	t.Log("1->2")
	result := lib.Compare(loadJSON(t, "json/prStateQuery1.json"), loadJSON(t, "json/prStateQuery2.json"))
	assert.Empty(t, result)

	t.Log("2->3")
	result = lib.Compare(loadJSON(t, "json/prStateQuery2.json"), loadJSON(t, "json/prStateQuery3.json"))
	assert.Empty(t, result)

	t.Log("3->4")
	result = lib.Compare(loadJSON(t, "json/prStateQuery3.json"), loadJSON(t, "json/prStateQuery4.json"))
	assert.ElementsMatch(t, []lib.Change{{lib.ALL_CHECKS_SUCCESS, "All status checks successfull. PR ready to go 🚀"}}, result)

	t.Log("4->5")
	result = lib.Compare(loadJSON(t, "json/prStateQuery4.json"), loadJSON(t, "json/prStateQuery5.json"))
	assert.ElementsMatch(t, []lib.Change{{lib.REVIEW_CHANGE, "joecommenter COMMENTED your PR"}}, result)

	t.Log("5->6")
	result = lib.Compare(loadJSON(t, "json/prStateQuery5.json"), loadJSON(t, "json/prStateQuery6.json"))
	assert.ElementsMatch(t, []lib.Change{{lib.REVIEW_CHANGE, "samapprover APPROVED your PR"}}, result)

	t.Log("6->7")
	result = lib.Compare(loadJSON(t, "json/prStateQuery6.json"), loadJSON(t, "json/prStateQuery7.json"))
	assert.ElementsMatch(t, []lib.Change{{lib.NEW_COMMIT, ""}}, result)

	t.Log("failure 1->2")
	result = lib.Compare(loadJSON(t, "json/prStateQueryfail1.json"), loadJSON(t, "json/prStateQueryfail2.json"))
	assert.ElementsMatch(t, []lib.Change{
		{lib.CHECK_FAILURE, "ci/qa FAILURE 🔥"},
		{lib.CHECK_FAILURE, "ci/tools/all FAILURE 🔥"},
		{lib.CHECK_FAILURE, "ci/tools/lint FAILURE 🔥"}},
		result,
	)

	t.Log("failure 2->3")
	result = lib.Compare(loadJSON(t, "json/prStateQueryfail2.json"), loadJSON(t, "json/prStateQueryfail3.json"))
	assert.Empty(t, result)
}

func loadJSON(t *testing.T, path string) *lib.RepositoryView {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var view lib.RepositoryView

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &view)
	require.NoError(t, err)

	return &view
}
