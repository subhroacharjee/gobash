package commands

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"unicode"
)

func (c Command) IsValidForCd() bool {
	return strings.HasPrefix(c.raw, "cd ") || c.raw == "cd"
}

func (c Command) Cd() (string, error) {
	args := c.args

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	if len(args) < 1 {
		if err := os.Chdir(usr.HomeDir); err != nil {
			return "\r", err
		}
		return "", nil
	}

	targetPath := strings.TrimFunc(args[0], unicode.IsSpace)

	if strings.HasPrefix(targetPath, "~") {
		homePath := os.Getenv("HOME")
		targetPath = filepath.Join(homePath, targetPath[1:])
	}
	if err := os.Chdir(targetPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("cd: %s: No such file or directory", targetPath)
		}

		return "", err
	}
	return "", nil
}

func (c Command) DescribeCd() string {
	return "cd is a shell builtin"
}
