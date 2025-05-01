package commands

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"unicode"
)

type Command struct {
	raw   string
	paths []string
}

func RunCommand(raw string) (string, error) {
	pathsStr := os.Getenv("PATH")
	paths := strings.Split(pathsStr, ":")
	cmd := Command{
		raw:   strings.TrimSuffix(raw, "\n"),
		paths: paths,
	}

	cmdName, err := cmd.searchFunctionToExecute()
	if err != nil {
		commandWithArgs := strings.Split(strings.TrimFunc(cmd.raw, unicode.IsSpace), " ")

		shellCmdWrapper := Command{raw: commandWithArgs[0], paths: paths}

		if _, err := shellCmdWrapper.searchCmdInPath(); err != nil {
			return "", fmt.Errorf("%s: command not found", strings.TrimSuffix(raw, "\n"))
		}

		shellCmd := exec.Command(commandWithArgs[0], commandWithArgs[1:]...)

		output, err := shellCmd.Output()
		if err != nil {
			return "", err
		} else {
			return string(output), nil
		}

	}

	v := reflect.ValueOf(cmd)
	result := v.MethodByName(cmdName).Call([]reflect.Value{})
	if !result[1].IsNil() {
		err = result[1].Interface().(error)
	}
	return result[0].String(), err
}

func (c Command) searchFunctionToExecute() (string, error) {
	v := reflect.ValueOf(c)
	t := reflect.TypeOf(c)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if strings.HasPrefix(method.Name, "IsValidFor") {
			result := v.MethodByName(method.Name).Call([]reflect.Value{})
			if len(result) > 0 && result[0].Bool() {
				cmdName := strings.Split(method.Name, "IsValidFor")[1]
				return cmdName, nil
			}

		}
	}

	return "", fmt.Errorf("not found")
}

func (c Command) searchCmdInPath() (string, error) {
	// fmt.Println("called")
	mainCmd := strings.Split(strings.TrimFunc(c.raw, unicode.IsSpace), " ")[0]
	for _, path := range c.paths {
		absPath := fmt.Sprintf("%s/%s", path, mainCmd)
		// fmt.Println(absPath)
		if _, err := os.Stat(absPath); err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				// TODO: figure it out
			}
		} else {
			return absPath, nil
		}

	}
	return "", fmt.Errorf("not found")
}
