package zql

import (
	"fmt"
	"reflect"
)

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
	for _, validator := range e.cfg.Validators {
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
	if counter > e.cfg.MaxPredicateNumber {
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
	for _, validator := range e.cfg.PredicateValidators {
		if validator != nil {
			err := validator.Validate(field, op, value, rt, rv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
