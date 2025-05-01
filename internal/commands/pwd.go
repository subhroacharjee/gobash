package commands

import (
	"fmt"
	"os"
	"strings"
)

func (c Command) IsValidForPwd() bool {
	return strings.HasPrefix(c.raw, "pwd ") || c.raw == "pwd"
}

func (c Command) Pwd() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n", pwd), nil
}

func (c Command) DescribePwd() string {
	return "pwd is a shell builtin"
}
