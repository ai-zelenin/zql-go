package zql

import (
	"fmt"
	"reflect"
)

type ValidatorOption func(i interface{}) error

type Validator interface {
	Validate(q *Query) error
}

type ExtendableValidator struct {
	cfg *ValidatorConfig
}

func NewExtendableValidator(cfg *ValidatorConfig) *ExtendableValidator {
	return &ExtendableValidator{cfg: cfg}
}

func (e *ExtendableValidator) Validate(q *Query) error {
	err := e.validateFilter(q)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExtendableValidator) validateFilter(q *Query) error {
	for _, predicate := range q.Filter {
		err := e.validatePredicate(predicate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ExtendableValidator) validatePredicate(p *Predicate) error {
	field := p.Field
	op := p.Op
	value := p.Value
	rt := reflect.TypeOf(value)
	valueType := TypeToString(rt)

	if e.cfg.CheckAcceptableFields {
		_, isFieldAcceptable := e.cfg.AcceptableFields[field]
		if !isFieldAcceptable {
			return NewError(fmt.Errorf("field %s is unacceptable", field), ErrCodeFieldUnacceptable)
		}
	}

	if e.cfg.CheckFieldValueTypeMap {
		isValueTypeGoodForField := e.cfg.FieldValueTypeMap[field] == valueType
		if !isValueTypeGoodForField {
			return NewError(fmt.Errorf("value type %v  unacceptable for field %s", valueType, field), ErrCodeValueTypeUnacceptableForField)
		}
	}

	if e.cfg.CheckOpValueTypeMap {
		isValueTypeGoodForOp := e.cfg.OpValueTypeMap[op] == valueType
		if !isValueTypeGoodForOp {
			return NewError(fmt.Errorf("value type %s unacceptable for op %s", valueType, op), ErrCodeValueTypeUnacceptableForOp)
		}
	}

	return nil
}
