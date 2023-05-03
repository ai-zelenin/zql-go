package compare

import (
	"reflect"
)

type PtrCaster struct {
	matrix map[reflect.Kind]CastFunc
}

func NewPtrCaster() *PtrCaster {
	var asPtr CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
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

	return &PtrCaster{
		matrix: map[reflect.Kind]CastFunc{
			reflect.Invalid:       asPtr,
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
			reflect.Chan:          asPtr,
			reflect.Func:          asPtr,
			reflect.Interface:     asPtr,
			reflect.Map:           asPtr,
			reflect.Pointer:       asPtr,
			reflect.Slice:         asPtr,
			reflect.String:        nil,
			reflect.Struct:        nil,
			reflect.UnsafePointer: nil,
		},
	}
}

func (c PtrCaster) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	toCmp := c.matrix[rv2.Kind()]
	if toCmp == nil {
		return nil, NewCompareError(rv1, rv2)
	}
	return toCmp(rv1, rv2)
}
