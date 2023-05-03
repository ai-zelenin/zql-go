package compare

import (
	"reflect"
)

type NilCaster struct {
	matrix map[reflect.Kind]CastFunc
}

func NewNilCaster() *NilCaster {
	var asNils CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		b1 := byte(0)
		b2 := byte(0)
		if !rv1.IsValid() || rv1.IsNil() {
			b1 = 1
		}
		if !rv2.IsValid() || rv2.IsNil() {
			b2 = 1
		}
		return NewOnlyEqGCV[byte](b1, b2), err
	}
	return &NilCaster{
		matrix: map[reflect.Kind]CastFunc{
			reflect.Invalid:    asNils,
			reflect.Bool:       nil,
			reflect.Int:        nil,
			reflect.Int8:       nil,
			reflect.Int16:      nil,
			reflect.Int32:      nil,
			reflect.Int64:      nil,
			reflect.Uint:       nil,
			reflect.Uint8:      nil,
			reflect.Uint16:     nil,
			reflect.Uint32:     nil,
			reflect.Uint64:     nil,
			reflect.Uintptr:    nil,
			reflect.Float32:    nil,
			reflect.Float64:    nil,
			reflect.Complex64:  nil,
			reflect.Complex128: nil,

			reflect.Array:     nil,
			reflect.Chan:      asNils,
			reflect.Func:      asNils,
			reflect.Interface: asNils,
			reflect.Map:       asNils,
			reflect.Pointer:   asNils,
			reflect.Slice:     asNils,

			reflect.String:        nil,
			reflect.Struct:        nil,
			reflect.UnsafePointer: nil,
		},
	}
}

func (c NilCaster) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	toCmp := c.matrix[rv2.Kind()]
	if toCmp == nil {
		return nil, NewCompareError(rv1, rv2)
	}
	return toCmp(rv1, rv2)
}
