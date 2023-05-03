package compare

import (
	"reflect"
	"time"
)

type StructCaster struct {
	matrix map[reflect.Kind]CastFunc
}

func NewStructCaster() *StructCaster {
	var timeType = reflect.TypeOf(time.Now())
	var asStructs CastFunc = func(rv1 reflect.Value, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
		rv1 = reflect.Indirect(rv1)
		rv2 = reflect.Indirect(rv2)
		t1 := rv1.Type()
		t2 := rv2.Type()
		if t1 != t2 {
			return nil, NewCompareError(rv1, rv2)
		}
		switch t1 {
		case timeType:
			cv1, ok := rv1.Interface().(time.Time)
			if !ok {
				return nil, NewCompareError(rv1, rv2)
			}
			cv2, ok := rv2.Interface().(time.Time)
			if !ok {
				return nil, NewCompareError(rv1, rv2)
			}
			return GCV[int64]{
				V1: cv1.UnixNano(),
				V2: cv2.UnixNano(),
			}, nil
		}
		return nil, NewCompareError(rv1, rv2)
	}
	return &StructCaster{
		matrix: map[reflect.Kind]CastFunc{
			reflect.Struct:  asStructs,
			reflect.Pointer: asStructs,
		},
	}
}

func (c StructCaster) AsComparableValues(rv1, rv2 reflect.Value) (cmpValues ComparableValues, err error) {
	toCmp := c.matrix[rv2.Kind()]
	if toCmp == nil {
		return nil, ErrTypesIsNotComparable
	}
	return toCmp(rv1, rv2)
}
