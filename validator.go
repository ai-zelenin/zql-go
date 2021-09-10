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
	validators          []Validator
	predicateValidators []ValidatorPredicate
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
		validators: []Validator{
			NewValidatorOrders(),
		},
		predicateValidators: []ValidatorPredicate{
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
	e.validators = append(e.validators, v)
}

func (e *ExtendableValidator) AddPredicateValidator(v ValidatorPredicate) {
	e.predicateValidators = append(e.predicateValidators, v)
}

func (e *ExtendableValidator) SetupValidatorForModel(model interface{}, tagName string) {
	modelDescription := ReflectModelDescription(model, tagName)
	for fieldName, valueType := range modelDescription {
		e.FieldsValidator.AddField(fieldName)
		e.ValuesValidator.AddFieldValuePair(fieldName, valueType)
	}
}

func (e *ExtendableValidator) Validate(q *Query) error {
	err := e.validateFilter(q)
	if err != nil {
		return err
	}
	for _, validator := range e.validators {
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
		_, err := predicate.Walk(func(parent Node, current Node, lvl int) (Node, error) {
			counter++
			currentPr := current.(*Predicate)
			err := e.validatePredicate(currentPr)
			if err != nil {
				return nil, err
			}
			return nil, nil
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
	for _, validator := range e.predicateValidators {
		if validator != nil {
			err := validator.Validate(field, op, value, rt, rv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
