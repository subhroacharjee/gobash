package commands

import (
	"fmt"
	"reflect"
	"strings"
)

type Command struct {
	raw string
}

func RunCommand(raw string) (string, error) {
	cmd := Command{
		raw: strings.TrimSuffix(raw, "\n"),
	}

	v := reflect.ValueOf(cmd)
	t := reflect.TypeOf(cmd)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if strings.HasPrefix(method.Name, "IsValidFor") {
			result := v.MethodByName(method.Name).Call([]reflect.Value{})
			if len(result) > 0 && result[0].Bool() {
				cmdName := strings.Split(method.Name, "IsValidFor")[1]
				result := v.MethodByName(cmdName).Call([]reflect.Value{})
				return result[0].String(), result[1].Interface().(error)
			}
		}
	}
	return "", fmt.Errorf("%s: command not found", strings.TrimSuffix(raw, "\n"))
}
