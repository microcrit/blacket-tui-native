package util

import (
	"strings"
)

func ParseCookie(cookie string) string {
	return strings.Split(cookie, ";")[0][6:]
}
