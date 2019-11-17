package shlex

import (
	"regexp"
	"strings"
)

var patFindUnsafe = regexp.MustCompile(`[^\w@%+=:,./-]`)

// copy from python stdlib shlex
func Quote(s string) string {
	if len(s) == 0 {
		return "''"
	}
	if patFindUnsafe.FindString(s) == "" {
		return s
	}
	return "'" + strings.Replace(s, "'", `'"'"'`, -1) + "'"
}
