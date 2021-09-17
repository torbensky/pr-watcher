package prwatcher_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	prwatcher "github.com/torbensky/pr-watcher"
)

func TestCompare(t *testing.T) {
	t.Log("1->2")
	result, err := prwatcher.Compare(loadJSON(t, "json/prStateQuery1.json"), loadJSON(t, "json/prStateQuery2.json"))
	require.NoError(t, err)
	assert.Equal(t, []prwatcher.ChangeType{prwatcher.NONE}, result)

	t.Log("2->3")
	result, err = prwatcher.Compare(loadJSON(t, "json/prStateQuery2.json"), loadJSON(t, "json/prStateQuery3.json"))
	require.NoError(t, err)
	assert.Equal(t, []prwatcher.ChangeType{prwatcher.CHECK_CHANGE, prwatcher.CHECK_CHANGE, prwatcher.CHECK_CHANGE}, result)

	t.Log("3->4")
	result, err = prwatcher.Compare(loadJSON(t, "json/prStateQuery3.json"), loadJSON(t, "json/prStateQuery4.json"))
	require.NoError(t, err)
	assert.Equal(t, []prwatcher.ChangeType{}, result)

	t.Log("4->5")
	result, err = prwatcher.Compare(loadJSON(t, "json/prStateQuery4.json"), loadJSON(t, "json/prStateQuery5.json"))
	require.NoError(t, err)
	assert.Equal(t, []prwatcher.ChangeType{prwatcher.REVIEW_CHANGE}, result)

	t.Log("5->6")
	result, err = prwatcher.Compare(loadJSON(t, "json/prStateQuery5.json"), loadJSON(t, "json/prStateQuery6.json"))
	require.NoError(t, err)
	assert.Equal(t, []prwatcher.ChangeType{prwatcher.REVIEW_CHANGE}, result)

	t.Log("6->7")
	result, err = prwatcher.Compare(loadJSON(t, "json/prStateQuery6.json"), loadJSON(t, "json/prStateQuery7.json"))
	require.NoError(t, err)
	assert.Equal(t, []prwatcher.ChangeType{prwatcher.NEW_COMMIT}, result)
}

func loadJSON(t *testing.T, path string) *prwatcher.RepositoryView {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var view prwatcher.RepositoryView

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)
	err = json.Unmarshal(byteValue, &view)
	require.NoError(t, err)

	return &view
}
