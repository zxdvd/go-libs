package utils

import "strings"

func StrReplace(tmpl string, data map[string]string) string {
	args := make([]string, 0, len(data)*2)
	for k, v := range data {
		k = "{" + k + "}"
		args = append(args, k, v)
	}
	return strings.NewReplacer(args...).Replace(tmpl)
}
