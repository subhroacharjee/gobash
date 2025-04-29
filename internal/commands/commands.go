package commands

import (
	"fmt"
	"strings"
)

type Command struct {
	raw string
}

func NewCommand(raw string) (*Command, error) {
	return nil, fmt.Errorf("%s: command not found", strings.TrimSuffix(raw, "\n"))
}
