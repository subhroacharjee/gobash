package commands

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func (c Command) IsValidForType() bool {
	return strings.HasPrefix(c.raw, "type ") || c.raw == "type"
}

func (c Command) Type() (string, error) {
	res := strings.Split(c.raw, "type ")
	if len(res) <= 1 {
		return "", fmt.Errorf("type needs atleast one value")
	}

	cmd := strings.TrimLeftFunc(res[1], unicode.IsSpace)
	// fmt.Println(cmd)
	cmdBase := Command{raw: cmd, paths: c.paths}
	cmdName, err := cmdBase.searchFunctionToExecute()
	if err != nil {
		absPath, err := cmdBase.searchCmdInPath()
		if err != nil {
			return "", fmt.Errorf("%s: not found", strings.TrimSuffix(cmd, "\n"))
		}

		return fmt.Sprintf("%s is %s", cmd, absPath), nil
	}
	askedCmdName := "Describe" + cmdName

	v := reflect.ValueOf(cmdBase)

	values := v.MethodByName(askedCmdName).Call([]reflect.Value{})
	if len(values) == 0 {
		return fmt.Sprintf("%s: doesnt have any type description", cmd), nil
	}
	return values[0].String(), nil
}

func (c Command) DescribeType() string {
	return "type is a shell builtin"
}
