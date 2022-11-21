package cmd_test

import (
	"ahkpm/src/cmd"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSearchResultsTable(t *testing.T) {
	items := []cmd.SearchResponseItem{
		{
			FullName:    "test/test",
			Description: "test",
		},
		{
			FullName:    "something/or-other",
			Description: "frobnicates the quux",
		},
	}

	var expectedBuilder strings.Builder
	expectedBuilder.WriteString("Name                         \tDescription         \n")
	expectedBuilder.WriteString("-----------------------------\t--------------------\n")
	expectedBuilder.WriteString("github.com/test/test         \ttest                \n")
	expectedBuilder.WriteString("github.com/something/or-other\tfrobnicates the quux\n")
	expected := expectedBuilder.String()

	actual := cmd.GetSearchResultsTable(items)

	assert.Equal(t, expected, actual)
}
