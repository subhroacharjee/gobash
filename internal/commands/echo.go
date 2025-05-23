package commands

import (
	"strings"
	"unicode"
)

func (c Command) IsValidForEcho() bool {
	return strings.HasPrefix(c.raw, "echo ") || c.raw == "echo"
}

func (c Command) Echo() (string, error) {
	data := strings.Join(c.args, "")
	return strings.TrimLeftFunc(data, unicode.IsSpace) + "\n", nil
}

func (c Command) DescribeEcho() string {
	return "echo is a shell builtin"
}
