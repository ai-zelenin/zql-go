package zql

import (
	"reflect"
	"strings"
)

func isNil(c interface{}) bool {
	v := reflect.ValueOf(c)
	return c == nil || (v.Kind() == reflect.Ptr && v.IsNil())
}

func FieldsFromModel(model interface{}) map[string]string {
	result := make(map[string]string, 0)
	rv := reflect.ValueOf(model)
	rt := rv.Type()
	for i := 0; i < rt.Len(); i++ {
		rf := rt.Field(i)
		result[rf.Name] = TypeToString(rf.Type)
	}
	return result
}

func TypeToString(t reflect.Type) string {
	return strings.Split(t.String(), ".")[1]
}
