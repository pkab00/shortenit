package stringutils

import "strings"

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
