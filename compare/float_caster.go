package compare

import (
	"reflect"
)

type FloatCaster struct {
	matrix map[reflect.Kind]CastFunc
}

func NewFloatCaster() *FloatCaster {
	var asFloat CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		return GCV[float64]{
			V1: rv1.Float(),
			V2: rv2.Float(),
		}, err
	}
	var float64Type = reflect.TypeOf(float64(1))
	var asNumber CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		return GCV[float64]{
			V1: rv1.Float(),
			V2: rv2.Convert(float64Type).Float(),
		}, err
	}
	return &FloatCaster{
		matrix: map[reflect.Kind]CastFunc{
			reflect.Int:     asNumber,
			reflect.Int8:    asNumber,
			reflect.Int16:   asNumber,
			reflect.Int32:   asNumber,
			reflect.Int64:   asNumber,
			reflect.Uint:    asNumber,
			reflect.Uint8:   asNumber,
			reflect.Uint16:  asNumber,
			reflect.Uint32:  asNumber,
			reflect.Uint64:  asNumber,
			reflect.Uintptr: asNumber,
			reflect.Float32: asFloat,
			reflect.Float64: asFloat,
		},
	}
}

func (c FloatCaster) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	toCmp := c.matrix[rv2.Kind()]
	if toCmp == nil {
		return nil, ErrTypesIsNotComparable
	}
	return toCmp(rv1, rv2)
}
