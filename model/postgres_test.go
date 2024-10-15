package model_test

import (
	"testing"

	"github.com/mattermost/mattermost-dbt/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPgIndexes(t *testing.T) {
	indexSlice := []*model.PgIndex{
		{IndexName: "index1"},
		{IndexName: "index2"},
		{IndexName: "index3"},
	}

	indexes := model.BuildPgIndexesFromSlice(indexSlice)
	require.Len(t, indexes, 3)

	indexNames := indexes.GetNames()
	assert.Contains(t, indexNames, "index1")
	assert.Contains(t, indexNames, "index2")
	assert.Contains(t, indexNames, "index3")
}
