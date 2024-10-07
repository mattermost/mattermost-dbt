package model_test

import (
	"testing"

	"github.com/mattermost/mattermost-dbt/model"
	"github.com/stretchr/testify/require"
)

func TestShortBuildHash(t *testing.T) {
	var testCases = []struct {
		testName string
		hash     string
		expected string
	}{
		{"empty", "", ""},
		{"long hash", "hash123456789", "hash123"},
		{"6 chars", "hash12", "hash12"},
		{"7 chars", "hash123", "hash123"},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			model.BuildHash = tc.hash
			require.Equal(t, tc.expected, model.ShortBuildHash())
		})
	}
}
