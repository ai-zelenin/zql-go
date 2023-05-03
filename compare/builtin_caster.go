package compare

import (
	"reflect"
)

type BuiltInCaster struct {
	matrix map[reflect.Kind]CastFunc
}

func NewBuiltInCaster() *BuiltInCaster {
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
	var asPtrEq CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		t1 := rv1.Type()
		t2 := rv2.Type()
		if t1 != t2 {
			return nil, NewCompareError(rv1, rv2)
		}
		b1 := uintptr(0)
		b2 := uintptr(0)
		if rv1.IsValid() && !rv1.IsNil() {
			b1 = rv1.Pointer()
		}
		if rv2.IsValid() && !rv2.IsNil() {
			b2 = rv2.Pointer()
		}
		return NewOnlyEqGCV[uintptr](b1, b2), err
	}

	return &BuiltInCaster{
		matrix: map[reflect.Kind]CastFunc{
			reflect.Invalid:       asNils,
			reflect.Bool:          nil,
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
			reflect.Chan:          asPtrEq,
			reflect.Func:          asPtrEq,
			reflect.Interface:     asPtrEq,
			reflect.Map:           asPtrEq,
			reflect.Pointer:       asPtrEq,
			reflect.Slice:         asPtrEq,
			reflect.String:        nil,
			reflect.Struct:        nil,
			reflect.UnsafePointer: nil,
		},
	}
}

func (c BuiltInCaster) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	toCmp := c.matrix[rv2.Kind()]
	if toCmp == nil {
		return nil, ErrTypesIsNotComparable
	}
	return toCmp(rv1, rv2)
}
