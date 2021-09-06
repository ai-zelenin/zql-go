package zql

import (
	"reflect"
)

func isNil(c interface{}) bool {
	v := reflect.ValueOf(c)
	return c == nil || (v.Kind() == reflect.Ptr && v.IsNil())
}
