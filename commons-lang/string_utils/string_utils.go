package string_utils

import "regexp"

var regex = regexp.MustCompile(`[[:punct:]]`)

func RemovePunctuation(input string) string {
	return regex.ReplaceAllString(input, "")
}
