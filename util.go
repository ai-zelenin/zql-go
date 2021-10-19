package zql

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type ValueType string

const (
	ValueTypeString = "string"
	ValueTypeNumber = "number"
	ValueTypeBool   = "bool"
)

var TimeReflectedType = reflect.TypeOf(time.Now())

func IsNilValue(c interface{}) bool {
	v := reflect.ValueOf(c)
	return c == nil || (v.Kind() == reflect.Ptr && v.IsNil())
}

func IsCompareOp(op string) bool {
	switch op {
	case EQ, NEQ, GT, GTE, LT, LTE:
		return true
	}
	return false
}

func ReflectModelDescription(model interface{}, tagName string) map[string]string {
	rv := reflect.ValueOf(model)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		rt = rt.Elem()
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	result := make(map[string]string, 0)
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		var name string
		if tagName == "" {
			name = rf.Name
		} else {
			tagValue := rf.Tag.Get(tagName)
			if tagValue != "" && tagValue != "-" {
				commaIdx := strings.Index(tagValue, ",")
				if commaIdx > 0 {
					name = tagValue[:commaIdx]
				} else {
					name = tagValue
				}
			}
		}
		result[name] = ValueTypeToString(rf.Type)
	}
	fmt.Println(result)
	return result
}

func ReflectValueTypeName(i interface{}) string {
	rt := reflect.TypeOf(i)
	return ValueTypeToString(rt)
}

func ValueTypeToString(rt reflect.Type) string {
	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ValueTypeNumber
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ValueTypeNumber
	case reflect.Float32, reflect.Float64:
		return ValueTypeNumber
	case reflect.Bool:
		return ValueTypeBool
	case reflect.String:
		return ValueTypeString
	}
	if rt == TimeReflectedType {
		return ValueTypeString
	}
	return rt.String()
}
