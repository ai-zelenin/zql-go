package zql

import (
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
	result := make(map[string]string, 0)
	rt := rv.Type()
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
				}
			}
		}
		result[name] = ValueTypeToString(rf.Type)
	}
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
