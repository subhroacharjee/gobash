package commands

import (
	"strings"
	"unicode"
)

func (c Command) IsValidForEcho() bool {
	return strings.HasPrefix(c.raw, "echo ") || c.raw == "echo"
}

func (c Command) Echo() (string, error) {
	data := strings.Split(c.raw, "echo")[1]
	return strings.TrimLeftFunc(data, unicode.IsSpace), nil
}
