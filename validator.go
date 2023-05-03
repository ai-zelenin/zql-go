package zql

import (
	"fmt"
	"reflect"
)

const DefaultMaxPredicateNumber = 20

type Validator interface {
	Validate(q *Query) error
}

type ExtendableValidator struct {
	Validators          []Validator
	PredicateValidators []ValidatorPredicate
	FieldsValidator     *ValidatorPredicateFields
	ValuesValidator     *ValidatorPredicateValues
	OpsValidator        *ValidatorPredicateOps
	OrdersValidator     *ValidatorOrders
	MaxPredicateNumber  int
}

func NewExtendableValidator() *ExtendableValidator {
	fieldsValidator := NewValidatorPredicateFields()
	valueValidator := NewValidatorPredicateValues()
	opsValidator := NewValidatorPredicateOps()
	return &ExtendableValidator{
		Validators: []Validator{
			NewValidatorOrders(),
		},
		PredicateValidators: []ValidatorPredicate{
			fieldsValidator,
			valueValidator,
			opsValidator,
		},
		FieldsValidator:    fieldsValidator,
		ValuesValidator:    valueValidator,
		OpsValidator:       opsValidator,
		MaxPredicateNumber: DefaultMaxPredicateNumber,
	}
}

func (e *ExtendableValidator) AddValidator(v Validator) {
	e.Validators = append(e.Validators, v)
}

func (e *ExtendableValidator) AddPredicateValidator(v ValidatorPredicate) {
	e.PredicateValidators = append(e.PredicateValidators, v)
}

func (e *ExtendableValidator) SetupValidatorForModel(model interface{}, tagName string, fieldDescFunc FieldDescFunc) {
	modelDescription := ReflectModelDescription(model, tagName, fieldDescFunc)
	for _, fieldDesc := range modelDescription {
		e.FieldsValidator.AddField(fieldDesc)
		e.ValuesValidator.AddFieldValuePair(fieldDesc)
	}
}

func (e *ExtendableValidator) Validate(q *Query) error {
	err := e.validateFilter(q)
	if err != nil {
		return err
	}
	for _, validator := range e.Validators {
		err = validator.Validate(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ExtendableValidator) validateFilter(q *Query) error {
	counter := 0
	for _, predicate := range q.Filter {
		counter++
		err := predicate.Walk(func(parent Node, current Node, lvl int) error {
			counter++
			currentPr := current.(*Predicate)
			err := e.validatePredicate(currentPr)
			if err != nil {
				return err
			}
			return nil
		}, nil, 0)
		if err != nil {
			return err
		}
	}
	if counter > e.MaxPredicateNumber {
		return NewError(fmt.Errorf("to many prediactes in filter"), ErrCodeTooManyPredicatesInFilter)
	}
	return nil
}

func (e *ExtendableValidator) validatePredicate(p *Predicate) error {
	field := p.Field
	op := p.Op
	value := p.Value
	rv := reflect.ValueOf(value)
	rt := reflect.TypeOf(value)
	for _, validator := range e.PredicateValidators {
		if validator != nil {
			err := validator.Validate(field, op, value, rt, rv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
