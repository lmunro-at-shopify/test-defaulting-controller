package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FormatLabel(t *testing.T) {
	testLabelName := "test-label"
	formatted := formatLabel(testLabelName)
	assert.Equal(t, testLabelName, formatted, "no change")

	separators := "/:"
	formatted = formatLabel(separators)
	assert.Equal(t, strings.Repeat("__", len(separators)), formatted, "separators should be replaced by __")

	invalidCharacters := "() <>$"
	formatted = formatLabel(invalidCharacters)
	assert.Equal(t, strings.Repeat("_", len(invalidCharacters)), formatted, "invalid characters are replaced")

	uppercaseString := strings.Repeat("A", 63)
	formatted = formatLabel(uppercaseString)
	assert.Equal(t, strings.Repeat("a", 63), formatted, "uppercase characters are downcased")
}

func Test_NormalizeLabels(t *testing.T) {
	tooLongString := strings.Repeat("a", maxLabelLength*maxLabelNum+1)
	_, err := NormalizeLabel(tooLongString)
	require.Error(t, err, "Labels over maxLabelNum*maxLabelLength should return an error")

	longLabel := strings.Repeat("a", maxLabelLength*maxLabelNum)
	labels, err := NormalizeLabel(longLabel)
	require.Nil(t, err, "There should be no error")
	require.Exactly(t, maxLabelNum, len(labels), fmt.Sprintf("There should be exactly %d label elements", maxLabelNum))
	require.Exactly(t, maxLabelLength, len(labels[0]), fmt.Sprintf("The first label should be exactly %d characters long", maxLabelLength))
	require.Exactly(t, maxLabelLength, len(labels[1]), fmt.Sprintf("The second label should be exactly %d characters long", maxLabelLength))

	shortLabel := strings.Repeat("a", 3)
	labels, err = NormalizeLabel(shortLabel)
	require.Nil(t, err, "There should be no error")
	require.Exactly(t, maxLabelNum, len(labels), fmt.Sprintf("There should be exactly %d label elements", maxLabelNum))
	require.Equal(t, shortLabel, labels[0], "The first element should be a short label")
	require.Equal(t, "", labels[1], "The second element should be an empty string")

	overflowLabel := strings.Repeat("a", maxLabelLength+5)
	labels, err = NormalizeLabel(overflowLabel)
	require.Nil(t, err, "There should be no error")
	require.Exactly(t, maxLabelNum, len(labels), fmt.Sprintf("There should be exactly %d label elements", maxLabelNum))
	require.Equal(t, strings.Repeat("a", maxLabelLength), labels[0], "The first element should be a full label")
	require.Equal(t, strings.Repeat("a", 5), labels[1], "The second element should be the remaining characters")
}
