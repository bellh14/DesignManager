package utils

import (
	"encoding/json"
	"os"
	"fmt"
	"reflect"
)

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func WriteBashVariable(file *os.File, name string, value any) {
	file.WriteString(fmt.Sprintf("%s=%v\n", name, value))
}

func WriteStructOfBashVariables(values reflect.Value, file *os.File) {
	for i := 0; i < values.NumField(); i++ {
		value := values.Field(i)
		name := values.Type().Field(i).Name
		WriteBashVariable(file, name, value.Interface())
	}
}