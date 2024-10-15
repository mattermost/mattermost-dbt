package main

import (
	"testing"

	"github.com/mattermost/mattermost-dbt/model"
	"github.com/stretchr/testify/assert"
)

func TestPgIndexes(t *testing.T) {
	sourceIndexes := model.BuildPgIndexesFromSlice([]*model.PgIndex{
		{IndexName: "index1"},
		{IndexName: "index2"},
		{IndexName: "index3"},
		{IndexName: "index5"},
	})
	assert.Len(t, sourceIndexes, 4)

	targetIndexes := model.BuildPgIndexesFromSlice([]*model.PgIndex{
		{IndexName: "index1"},
		{IndexName: "index2"},
		{IndexName: "index4"},
		{IndexName: "index5"},
		{IndexName: "index6"},
	})
	assert.Len(t, targetIndexes, 5)

	missing := compareIndexes(sourceIndexes, targetIndexes)
	assert.Len(t, missing, 1)
	assert.Contains(t, missing, "index3")

	missing = compareIndexes(targetIndexes, sourceIndexes)
	assert.Len(t, missing, 2)
	assert.Contains(t, missing, "index4")
	assert.Contains(t, missing, "index6")
}
