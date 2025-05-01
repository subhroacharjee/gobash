package commands

import (
	"os"
	"strings"
)

func (c Command) IsValidForPwd() bool {
	return strings.HasPrefix(c.raw, "pwd ") || c.raw == "pwd"
}

func (c Command) Pwd() (string, error) {
	return os.Getwd()
}

func (c Command) DescribePwd() string {
	return "pwd is a shell builtin"
}
