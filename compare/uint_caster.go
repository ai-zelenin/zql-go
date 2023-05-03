package compare

import (
	"reflect"
)

type UintCaster struct {
	matrix map[reflect.Kind]CastFunc
}

func NewUintCaster() *UintCaster {
	var asUInt CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		return GCV[uint64]{
			V1: rv1.Uint(),
			V2: rv2.Uint(),
		}, err
	}
	var float64Type = reflect.TypeOf(float64(1))
	var asNumber CastFunc = func(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		return GCV[float64]{
			V1: float64(rv1.Uint()),
			V2: rv2.Convert(float64Type).Float(),
		}, err
	}
	return &UintCaster{
		matrix: map[reflect.Kind]CastFunc{
			reflect.Int:     asNumber,
			reflect.Int8:    asNumber,
			reflect.Int16:   asNumber,
			reflect.Int32:   asNumber,
			reflect.Int64:   asNumber,
			reflect.Uint:    asUInt,
			reflect.Uint8:   asUInt,
			reflect.Uint16:  asUInt,
			reflect.Uint32:  asUInt,
			reflect.Uint64:  asUInt,
			reflect.Uintptr: asUInt,
			reflect.Float32: asNumber,
			reflect.Float64: asNumber,
		},
	}
}

func (c UintCaster) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	toCmp := c.matrix[rv2.Kind()]
	if toCmp == nil {
		return nil, ErrTypesIsNotComparable
	}
	return toCmp(rv1, rv2)
}
