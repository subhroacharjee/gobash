package commands

import (
	"os"
	"strconv"
	"strings"
)

func (c Command) IsValidForExit() bool {
	return strings.HasPrefix(c.raw, "exit ") || c.raw == "exit"
}

func (c Command) Exit() (string, error) {
	res := c.args
	if len(res) <= 1 {
		os.Exit(0)
	}
	stsCode := res[1]
	if len(stsCode) == 0 {
		os.Exit(0)
	}

	code, err := strconv.Atoi(stsCode)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(code)
	return "", nil
}

func (c Command) DescribeExit() string {
	return "exit is a shell builtin"
}
