package agent_adapter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func CatchPanic(f func()) (panicked bool, reseaon interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			reseaon = r
		}
	}()
	f()
	return
}

func CallWithJSON(function any, rawjson string) ([]reflect.Value, error) {
	funcValue := reflect.ValueOf(function)
	paramType := funcValue.Type().In(0)
	argPtr := reflect.New(paramType)
	if err := json.Unmarshal([]byte(rawjson), argPtr.Interface()); err != nil {
		return nil, err
	}

	var values []reflect.Value
	panicked, reseaon := CatchPanic(func() {
		values = funcValue.Call([]reflect.Value{argPtr.Elem()})
	})

	if panicked {
		return nil, fmt.Errorf("Panic: %v", reseaon)
	}

	return values, nil
}

type Parameter struct {
	Type        string
	Description string
}

func GetFunctionName(function any) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	return parts[len(parts)-1]
}

func GetFunctionDesc(function any) string {
	functionName := GetFunctionName(function)
	functionDesc := ""

	funcValue := reflect.ValueOf(function)
	paramType := funcValue.Type().In(0)

	parameters := make(map[string]Parameter)
	for i := 0; i < paramType.NumField(); i++ {
		field := paramType.Field(i)
		if field.Name == "Desc" {
			functionDesc = field.Tag.Get("desc")
			continue
		}
		parameters[field.Name] = Parameter{
			Type:        field.Type.String(),
			Description: field.Tag.Get("desc"),
		}
	}

	parameterNames := make([]string, 0, len(parameters))
	for name := range parameters {
		parameterNames = append(parameterNames, name)
	}

	stringParameters, err := json.Marshal(parameters)
	stringParameterNames, err := json.Marshal(parameterNames)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`{
		"type": "function",
		"function": {
			"name": "%s",
			"description": "%s",
			"parameters": {
				"type": "object",
				"properties": %s,
				"required": %s
			}
		}
	}`, functionName, functionDesc, stringParameters, stringParameterNames)
}
