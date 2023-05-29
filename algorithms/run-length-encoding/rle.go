package run_length_encoding

import (
	"fmt"
	"strconv"
	"strings"
)

func Encode(s string) string {
	if len(s) == 0 {
		return ""
	}
	sb := strings.Builder{}
	runLength := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			runLength++
		} else {
			sb.WriteString(fmt.Sprintf("%d%c", runLength, s[i-1]))
			runLength = 1
		}
	}
	if runLength > 0 {
		sb.WriteString(fmt.Sprintf("%d%c", runLength, s[len(s)-1]))
	}
	return sb.String()
}
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func Decode(s string) string {
	if len(s) < 2 {
		return ""
	}
	sb := strings.Builder{}
	i := 0
	for i <= len(s)-2 {
		j := i
		for j <= len(s)-1 && isDigit(s[j]) {
			j++
		}
		// j is the first non-digit
		str := s[i:j]
		n, _ := strconv.Atoi(str)
		for n > 0 {
			sb.WriteByte(s[j])
			n--
		}
		i = j + 1
	}
	return sb.String()
}
