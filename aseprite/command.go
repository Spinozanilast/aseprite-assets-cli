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

func CreateArgsFromStruct(s interface{}) []string {
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
			isUiEnabled := field.Interface()
			if !(isUiEnabled).(bool) {
				args = append(args, "-b")
			}
			continue
		}

		var key string
		if fieldType.Tag.Get("script") == "" {
			key = strings.ToLower(fieldType.Name)
		} else {
			key = fieldType.Tag.Get("script")
		}

		value := field.Interface()

		if fieldType.Tag.Get("format") == "quotes" {
			value = fmt.Sprintf("\"%v\"", value)
		}

		args = append(args, createScriptArgs(key, value)...)
	}

	return args
}

func createScriptArgs(key string, value any) []string {
	return []string{AsepriteScriptParamArg, fmt.Sprintf("%s=%v", key, value)}
}
