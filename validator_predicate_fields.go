package zql

import (
	"fmt"
	"reflect"
)

type ValidatorPredicateFields struct {
	Fields           map[string]FieldDesc
	CustomValidators []ValidatorPredicate
}

func NewValidatorPredicateFields() *ValidatorPredicateFields {
	return &ValidatorPredicateFields{
		Fields:           map[string]FieldDesc{},
		CustomValidators: []ValidatorPredicate{},
	}
}

func (e *ValidatorPredicateFields) AddField(f FieldDesc) {
	e.Fields[f.Name] = f
	if f.ValidateFunc != nil {
		e.CustomValidators = append(e.CustomValidators, f.ValidateFunc)
	}
}

func (e *ValidatorPredicateFields) Validate(field, op string, value interface{}, rt reflect.Type, rv reflect.Value) error {
	if field != "" {
		_, isFieldAcceptable := e.Fields[field]
		if isFieldAcceptable {
			return nil
		}
		for _, f := range e.CustomValidators {
			err := f.Validate(field, op, value, rt, rv)
			if err != nil {
				return err
			}
		}
		return NewError(fmt.Errorf("field %s is unacceptable", field), ErrCodeFieldUnacceptable)
	}
	return nil
}
