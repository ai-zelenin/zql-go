package zql

import (
	"reflect"
)

type ValidatorPredicate interface {
	Validate(field, op string, value interface{}, rt reflect.Type, rv reflect.Value) error
}

type ValidatePredicateFunc func(field, op string, value interface{}, rt reflect.Type, rv reflect.Value) error

func (v ValidatePredicateFunc) Validate(field, op string, value interface{}, rt reflect.Type, rv reflect.Value) error {
	return v(field, op, value, rt, rv)
}
