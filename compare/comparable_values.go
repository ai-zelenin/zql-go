package compare

type Comparable interface {
	EqInterface
	NeqInterface
	GtInterface
	GteInterface
	LtInterface
	LteInterface
}

type EqInterface interface {
	Eq(value any) (result bool, err error)
}

type NeqInterface interface {
	Neq(value any) (result bool, err error)
}

type GtInterface interface {
	Gt(value any) (result bool, err error)
}

type GteInterface interface {
	Gte(value any) (result bool, err error)
}

type LtInterface interface {
	Lt(value any) (result bool, err error)
}

type LteInterface interface {
	Lte(value any) (result bool, err error)
}

type InterfaceComparableValues struct {
	V1 any
	V2 any
}

func (i InterfaceComparableValues) Eq() (result bool, err error) {
	cmp, ok := i.V1.(EqInterface)
	if !ok {
		return false, ErrCompareOperationIsNotAcceptable
	}
	return cmp.Eq(i.V2)
}

func (i InterfaceComparableValues) Neq() (result bool, err error) {
	cmp, ok := i.V1.(NeqInterface)
	if !ok {
		return false, ErrCompareOperationIsNotAcceptable
	}
	return cmp.Neq(i.V2)
}

func (i InterfaceComparableValues) Gt() (result bool, err error) {
	cmp, ok := i.V1.(GtInterface)
	if !ok {
		return false, ErrCompareOperationIsNotAcceptable
	}
	return cmp.Gt(i.V2)
}

func (i InterfaceComparableValues) Gte() (result bool, err error) {
	cmp, ok := i.V1.(GteInterface)
	if !ok {
		return false, ErrCompareOperationIsNotAcceptable
	}
	return cmp.Gte(i.V2)
}

func (i InterfaceComparableValues) Lt() (result bool, err error) {
	cmp, ok := i.V1.(LtInterface)
	if !ok {
		return false, ErrCompareOperationIsNotAcceptable
	}
	return cmp.Lt(i.V2)
}

func (i InterfaceComparableValues) Lte() (result bool, err error) {
	cmp, ok := i.V1.(LteInterface)
	if !ok {
		return false, ErrCompareOperationIsNotAcceptable
	}
	return cmp.Lte(i.V2)
}
