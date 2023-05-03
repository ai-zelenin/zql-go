package compare

import (
	"reflect"
)

type IntCaster struct {
	matrix map[reflect.Kind]CastFunc
}

func NewIntCaster() *IntCaster {
	var asInt CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		return GCV[int64]{
			V1: rv1.Int(),
			V2: rv2.Int(),
		}, err
	}
	var float64Type = reflect.TypeOf(float64(1))
	var asNumber CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		return GCV[float64]{
			V1: float64(rv1.Int()),
			V2: rv2.Convert(float64Type).Float(),
		}, err
	}

	return &IntCaster{
		matrix: map[reflect.Kind]CastFunc{
			reflect.Invalid:       nil,
			reflect.Bool:          nil,
			reflect.Int:           asInt,
			reflect.Int8:          asInt,
			reflect.Int16:         asInt,
			reflect.Int32:         asInt,
			reflect.Int64:         asInt,
			reflect.Uint:          asNumber,
			reflect.Uint8:         asNumber,
			reflect.Uint16:        asNumber,
			reflect.Uint32:        asNumber,
			reflect.Uint64:        asNumber,
			reflect.Uintptr:       asNumber,
			reflect.Float32:       asNumber,
			reflect.Float64:       asNumber,
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

func (c IntCaster) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	toCmp := c.matrix[rv2.Kind()]
	if toCmp == nil {
		return nil, NewCompareError(rv1, rv2)
	}
	return toCmp(rv1, rv2)
}
