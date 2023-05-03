package compare

import (
	"fmt"
	"reflect"
)

type Caster interface {
	AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error)
}

type CastFunc func(rv1 reflect.Value, rv2 reflect.Value) (cmpValues ComparableValues, err error)

func (c CastFunc) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	return c(rv1, rv2)
}

type CommonCaster struct {
	matrix map[reflect.Kind]Caster
}

func NewCaster() *CommonCaster {
	intCaster := NewIntCaster()
	uintCaster := NewUintCaster()
	floatCaster := NewFloatCaster()
	boolCaster := NewBoolCaster()
	structCaster := NewStructCaster()
	nilCaster := NewNilCaster()
	ptrCaster := NewPtrCaster()
	buildInCaster := NewBuiltInCaster()
	var asString CastFunc = func(rv1 reflect.Value, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		return GCV[string]{
			V1: rv1.String(),
			V2: rv2.String(),
		}, nil
	}
	return &CommonCaster{
		matrix: map[reflect.Kind]Caster{
			reflect.Invalid: nilCaster,

			reflect.Bool: boolCaster,

			reflect.Int:   intCaster,
			reflect.Int8:  intCaster,
			reflect.Int16: intCaster,
			reflect.Int32: intCaster,
			reflect.Int64: intCaster,

			reflect.Uint:    uintCaster,
			reflect.Uint8:   uintCaster,
			reflect.Uint16:  uintCaster,
			reflect.Uint32:  uintCaster,
			reflect.Uint64:  uintCaster,
			reflect.Uintptr: uintCaster,

			reflect.Float32: floatCaster,
			reflect.Float64: floatCaster,

			reflect.Complex64:  nil,
			reflect.Complex128: nil,

			reflect.Array:         nil,
			reflect.Chan:          buildInCaster,
			reflect.Func:          buildInCaster,
			reflect.Interface:     buildInCaster,
			reflect.Map:           buildInCaster,
			reflect.Pointer:       ptrCaster,
			reflect.Slice:         buildInCaster,
			reflect.String:        asString,
			reflect.Struct:        structCaster,
			reflect.UnsafePointer: nil,
		},
	}
}

func (c *CommonCaster) AsComparableValues(v1, v2 any) (cmpValues ComparableValues, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in AsComparableValues func %v", r)
		}
	}()
	rv1 := reflect.ValueOf(v1)
	rv2 := reflect.ValueOf(v2)
	rvk2 := rv1.Kind()
	comparator := c.matrix[rvk2]
	if comparator == nil {
		return nil, NewCompareError(rv1, rv2)
	}
	return comparator.AsComparableValues(rv1, rv2)
}
