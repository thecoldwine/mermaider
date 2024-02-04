package internal

import (
	"strings"
)

func PostgresFlavorer(s string) string {
	return strings.ReplaceAll(s, "@schema", "$1")
}
