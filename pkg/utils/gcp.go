package utils

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	maxLabelLength = 63
	maxLabelNum    = 2
)

var (
	separatorChars = regexp.MustCompile("[/:]")
	forbiddenChars = regexp.MustCompile("[^0-9a-z._-]") // the ^ negates the rest of the characters
)

// NormalizeLabel enforces GCP limitations on labels and returns an array (length == maxLabelNum)
// of normalized strings. Only the first element of the array is garanteed to contain
// more than an empty string. The others are only filled in if the string passed as
// argument is long enough.
func NormalizeLabel(s string) ([maxLabelNum]string, error) {
	s = formatLabel(s)
	sLength := len([]rune(s))
	if sLength > maxLabelLength*maxLabelNum {
		return [maxLabelNum]string{}, fmt.Errorf("label overflow, length %d, maximum %d", sLength, maxLabelLength*maxLabelNum)
	}

	var normalizedList [maxLabelNum]string
	remaining := sLength
	start, end := 0, 0
	for i := range normalizedList {
		end += min(maxLabelLength, remaining)
		normalizedList[i] = s[start:end]

		remaining -= (end - start)
		if remaining == 0 {
			break
		}
		start += maxLabelLength
	}

	return normalizedList, nil
}

// Format strings to use as a label value.
// Labels values must be lowercase and 63 characters or less.
// replace all separator characters with double underscores, all other forbiddenChars with a single underscore.
// Cf. https://cloud.google.com/compute/docs/labeling-resources#requirements
func formatLabel(s string) string {
	s = strings.ToLower(s)
	s = separatorChars.ReplaceAllString(s, "__")

	return forbiddenChars.ReplaceAllString(s, "_")
}

// min is a utility funtion to compute the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
