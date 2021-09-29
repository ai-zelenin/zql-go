package zql

import (
	"fmt"
	"reflect"
)

type ValidatorPredicateOps struct {
	UnacceptableOps       map[string]bool
	RequireOpValueTypeMap map[string]string
	RequireOpValueKindMap map[string][]reflect.Kind
}

func NewValidatorPredicateOps(unacceptableOps ...string) *ValidatorPredicateOps {
	m := make(map[string]bool)
	for _, field := range unacceptableOps {
		m[field] = true
	}
	return &ValidatorPredicateOps{
		UnacceptableOps: m,
		RequireOpValueKindMap: map[string][]reflect.Kind{
			IN:    {reflect.Slice, reflect.Array},
			LIKE:  {reflect.String},
			ILIKE: {reflect.String},
		},
		RequireOpValueTypeMap: map[string]string{
			AND: ReflectValueTypeName(make([]*Predicate, 0)),
			OR:  ReflectValueTypeName(make([]*Predicate, 0)),
		},
	}
}

func (e *ValidatorPredicateOps) AddUnacceptableOp(op string) {
	e.UnacceptableOps[op] = true
}

func (e *ValidatorPredicateOps) Validate(field, op string, value interface{}, rt reflect.Type, rv reflect.Value) error {
	isOpUnacceptable := e.UnacceptableOps[op]
	if isOpUnacceptable {
		return NewError(fmt.Errorf("op %s is unacceptable", op), ErrCodeUnacceptableOp)
	}

	kinds, ok := e.RequireOpValueKindMap[op]
	if ok {
		var isValueKindGoodForOp bool
		for _, kind := range kinds {
			if rt.Kind() == kind {
				isValueKindGoodForOp = true
			}
		}
		if !isValueKindGoodForOp {
			return NewError(fmt.Errorf("value kind %s unacceptable for op %s", rt.Kind(), op), ErrCodeValueKindUnacceptableForOp)
		}
	}

	if rt == nil {
		return nil
	}
	valueType := ValueTypeToString(rt)
	requiredValueType, ok := e.RequireOpValueTypeMap[op]
	if ok {
		isValueTypeGoodForOp := requiredValueType == valueType
		if !isValueTypeGoodForOp {
			return NewError(fmt.Errorf("value type %s unacceptable for op %s", valueType, op), ErrCodeValueTypeUnacceptableForOp)
		}
	}
	return nil
}
