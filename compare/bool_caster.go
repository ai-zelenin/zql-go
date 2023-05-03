package compare

import (
	"reflect"
)

type BoolCaster struct {
	matrix map[reflect.Kind]CastFunc
}

func NewBoolCaster() *BoolCaster {
	var asBool CastFunc = func(rv1 reflect.Value, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		b1 := byte(0)
		b2 := byte(0)
		if rv1.Bool() {
			b1 = 1
		}
		if rv2.Bool() {
			b2 = 1
		}
		return NewGCV[byte](b1, b2), nil
	}
	return &BoolCaster{
		matrix: map[reflect.Kind]CastFunc{
			reflect.Invalid:       nil,
			reflect.Bool:          asBool,
			reflect.Int:           nil,
			reflect.Int8:          nil,
			reflect.Int16:         nil,
			reflect.Int32:         nil,
			reflect.Int64:         nil,
			reflect.Uint:          nil,
			reflect.Uint8:         nil,
			reflect.Uint16:        nil,
			reflect.Uint32:        nil,
			reflect.Uint64:        nil,
			reflect.Uintptr:       nil,
			reflect.Float32:       nil,
			reflect.Float64:       nil,
			reflect.Complex64:     nil,
			reflect.Complex128:    nil,
			reflect.Array:         nil,
			reflect.Chan:          nil,
			reflect.Func:          nil,
			reflect.Interface:     nil,
			reflect.Map:           nil,
			reflect.Pointer:       nil,
			reflect.Slice:         nil,
			reflect.String:        nil,
			reflect.Struct:        nil,
			reflect.UnsafePointer: nil,
		},
	}
}

func (c BoolCaster) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	toCmp := c.matrix[rv2.Kind()]
	if toCmp == nil {
		return nil, ErrTypesIsNotComparable
	}
	return toCmp(rv1, rv2)
}
