package zql

import (
	"fmt"
	"reflect"
)

type ValidatorPredicateFields struct {
	AcceptableFields map[string]bool
}

func NewValidatorPredicateFields(acceptableFields ...string) *ValidatorPredicateFields {
	m := make(map[string]bool)
	for _, field := range acceptableFields {
		m[field] = true
	}
	return &ValidatorPredicateFields{
		AcceptableFields: m,
	}
}

func (e *ValidatorPredicateFields) AddField(f string) {
	e.AcceptableFields[f] = true
}

func (e *ValidatorPredicateFields) Validate(field, _ string, _ interface{}, _ reflect.Type, _ reflect.Value) error {
	if field != "" {
		_, isFieldAcceptable := e.AcceptableFields[field]
		if !isFieldAcceptable {
			return NewError(fmt.Errorf("field %s is unacceptable", field), ErrCodeFieldUnacceptable)
		}
	}
	return nil
}
