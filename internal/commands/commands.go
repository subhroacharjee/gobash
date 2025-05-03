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
	raw     string
	paths   []string
	args    []string
	rawArgs string
}

func RunCommand(raw string) (string, error) {
	pathsStr := os.Getenv("PATH")
	paths := strings.Split(pathsStr, ":")
	commandWithArgs := strings.Split(strings.TrimFunc(strings.TrimSuffix(raw, "\n"), unicode.IsSpace), " ")
	if len(commandWithArgs[0]) == 0 {
		return "", nil
	}
	cmd := Command{
		raw:     commandWithArgs[0],
		paths:   paths,
		args:    commandWithArgs[1:],
		rawArgs: strings.Split(strings.TrimFunc(strings.TrimSuffix(raw, "\n"), unicode.IsSpace), commandWithArgs[0])[1],
	}

	if err := cmd.ParseArgs(); err != nil {
		return "", err
	}

	cmdName, err := cmd.searchFunctionToExecute()
	if err != nil {

		if _, err := cmd.searchCmdInPath(); err != nil {
			return "", fmt.Errorf("%s: command not found", strings.TrimSuffix(raw, "\n"))
		}

		shellCmd := exec.Command(cmd.raw, cmd.args...)
		shellCmd.Stdout = os.Stdout
		shellCmd.Stdin = os.Stderr

		err := shellCmd.Run()
		if err != nil {
			return "", nil
		} else {
			return "", nil
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

func (c *Command) ParseArgs() error {
	finalArgs := make([]string, 0)
	strArgs := strings.TrimFunc(strings.Join(c.args, " "), unicode.IsSpace)
	if len(strArgs) == 0 {
		c.args = finalArgs
		return nil
	}
	// fmt.Println(strArgs)

	arg := ""
	for i := 0; i < len(strArgs); i++ {
		r := strArgs[i]
		if r == '\'' {
			for i = i + 1; i < len(strArgs) && strArgs[i] != '\''; i++ {
				arg += string(strArgs[i])
			}
			if i == len(strArgs) && strArgs[i-1] != '\'' {
				return fmt.Errorf("unterminated single quotation")
			}
			finalArgs = append(finalArgs, arg)
			arg = ""
			continue
		} else if r == '"' {
			for i = i + 1; i < len(strArgs) && strArgs[i] != '"'; i++ {
				arg += string(strArgs[i])

				// if strArgs[i] == '\\' {
				// 	if i == len(strArgs)-1 {
				// 		return fmt.Errorf("unterminated double quotation")
				// 	}
				// 	arg += string(strArgs[i : i+1])
				// 	i++
				// } else {
				//
				// }
			}

			finalArgs = append(finalArgs, arg)
			arg = ""
			continue
		} else if string(r) == "\\" {
			if i == len(strArgs)-1 {
				return fmt.Errorf("unterminated double quotation")
			}
			arg += string(strArgs[i+1])
			i++
		} else if unicode.IsSpace(rune(r)) {
			if len(arg) != 0 {
				finalArgs = append(finalArgs, arg, " ")
				arg = ""
				continue
			}

			if !unicode.IsSpace(rune(finalArgs[len(finalArgs)-1][0])) {
				finalArgs = append(finalArgs, " ")
			}
		} else {
			arg += string(r)
		}
		// fmt.Println(">>>>>>>>>>>>>>>", arg, len(arg))
	}
	finalArgs = append(finalArgs, arg)

	c.args = finalArgs

	return nil
}
