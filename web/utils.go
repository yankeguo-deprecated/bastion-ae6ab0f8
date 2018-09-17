package web

import "strings"

func IsFormValueTrue(v string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(v)), "t")
}
