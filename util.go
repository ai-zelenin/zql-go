package zql

import (
	"reflect"
	"regexp"
	"strings"
	"time"
)

type ValueType string

const (
	ValueTypeString = "string"
	ValueTypeNumber = "number"
	ValueTypeBool   = "bool"
	ValueTypeStruct = "struct"
)

var TimeReflectedType = reflect.TypeOf(time.Now())

var SanitizeRegexp = regexp.MustCompile(`[^\w\d.]`)

func Sanitize(str string) string {
	return SanitizeRegexp.ReplaceAllString(str, "")
}

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

type FieldDescFunc func(desc FieldDesc, rf reflect.StructField) FieldDesc

type FieldDesc struct {
	Name         string
	Type         string
	ValidateFunc ValidatePredicateFunc
}

func ReflectModelDescription(model interface{}, tagName string, fieldDescFunc FieldDescFunc) []FieldDesc {
	rt := reflect.TypeOf(model)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		rt = rt.Elem()
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	result := make([]FieldDesc, 0)
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
		desc := FieldDesc{
			Name: name,
			Type: ValueTypeToString(rf.Type),
		}
		if fieldDescFunc != nil {
			desc = fieldDescFunc(desc, rf)
		}
		result = append(result, desc)
	}
	return result
}

func ReflectValueTypeName(i interface{}) string {
	rt := reflect.TypeOf(i)
	return ValueTypeToString(rt)
}

func ValueTypeToString(rt reflect.Type) string {
	if rt == TimeReflectedType {
		return ValueTypeString
	}
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
	case reflect.Struct, reflect.Map:
		return ValueTypeStruct
	}
	return rt.String()
}
