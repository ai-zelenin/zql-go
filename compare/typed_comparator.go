package compare

import (
	"fmt"
)

const (
	Eq  = "eq"
	Neq = "neq"
	Gt  = "gt"
	Gte = "gte"
	Lt  = "lt"
	Lte = "lte"
)

type TypedComparator struct {
	*CommonCaster
}

func NewTypedComparator() *TypedComparator {
	return &TypedComparator{
		CommonCaster: NewCaster(),
	}
}

func (c *TypedComparator) Eq(v1, v2 any) (result bool, err error) {
	cmp, err := c.AsComparableValues(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp.Eq()
}

func (c *TypedComparator) Neq(v1, v2 any) (result bool, err error) {
	cmp, err := c.AsComparableValues(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp.Neq()
}

func (c *TypedComparator) Gt(v1, v2 any) (result bool, err error) {
	cmp, err := c.AsComparableValues(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp.Gt()
}

func (c *TypedComparator) Gte(v1, v2 any) (result bool, err error) {
	cmp, err := c.AsComparableValues(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp.Gte()
}

func (c *TypedComparator) Lt(v1, v2 any) (result bool, err error) {
	cmp, err := c.AsComparableValues(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp.Lt()
}

func (c *TypedComparator) Lte(v1, v2 any) (result bool, err error) {
	cmp, err := c.AsComparableValues(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp.Lte()
}

func (c *TypedComparator) Op(op string, v1, v2 any) (result bool, err error) {
	switch op {
	case Eq:
		return c.Eq(v1, v2)
	case Neq:
		return c.Neq(v1, v2)
	case Gt:
		return c.Gt(v1, v2)
	case Gte:
		return c.Gte(v1, v2)
	case Lt:
		return c.Lt(v1, v2)
	case Lte:
		return c.Lte(v1, v2)
	}
	return false, fmt.Errorf("no operation with name %s", op)
}
