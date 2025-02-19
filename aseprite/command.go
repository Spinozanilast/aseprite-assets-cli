package aseprite

import (
	"fmt"
	"reflect"
	"strings"
)

const AsepriteScriptParamArg = "--script-param"

type AsepriteCommand interface {
	GetArgs() []string
	GetScriptName() string
}

func createArgsFromStruct(s interface{}) []string {
	args := []string{}
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if fieldType.Tag.Get("script") == "ignore" {
			continue
		}
		if fieldType.Name == "Ui" {
			// Ui flag might be a standalone flag.
			args = append(args, "-b")
			continue
		}

		var key string
		if fieldType.Tag.Get("script") == "" {
			key = strings.ToLower(fieldType.Name)
		} else {
			key = fieldType.Tag.Get("script")
		}

		value := field.Interface()

		args = append(args, CreateScriptArgs(key, value)...)
	}

	return args
}

func CreateScriptArgs(key string, value any) []string {
	return []string{AsepriteScriptParamArg, fmt.Sprintf("%s=%v", key, value)}
}
