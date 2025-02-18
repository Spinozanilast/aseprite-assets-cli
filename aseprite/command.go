package aseprite

import (
	"fmt"
	"reflect"
	"strings"
)

const AsepriteScriptParamArg = "-script-param"

type AsepriteCommand interface {
	GetArgs() []string
}

func createArgsFromStruct(s interface{}) []string {
	args := []string{}
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if fieldType.Tag.Get("aseprite") == "ignore" {
			continue
		}
		if fieldType.Name == "Ui" {
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

		args = append(args, CreateScriptArg(key, value))
	}

	return args
}

func CreateScriptArg(key string, value interface{}) string {
	return fmt.Sprintf("%s %s=%v", AsepriteScriptParamArg, key, value)
}
