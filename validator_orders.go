package zql

import (
	"fmt"
)

type ValidatorOrders struct {
}

func NewValidatorOrders() *ValidatorOrders {
	return &ValidatorOrders{}
}

func (v *ValidatorOrders) Validate(q *Query) error {
	for _, order := range q.Orders {
		err := v.validateOrder(order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *ValidatorOrders) validateOrder(o *Order) error {
	switch o.Direction {
	case ASC, DESC:
	default:
		return NewError(fmt.Errorf("unacceptable order direction"), ErrCodeUnacceptableOrderDirection)
	}
	return nil
}
