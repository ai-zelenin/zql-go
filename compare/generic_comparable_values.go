package compare

type ComparableValueTypes interface {
	int | int8 | int16 | int32 | int64 | uint | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string | byte
}

type ComparableValues interface {
	Eq() (result bool, err error)
	Neq() (result bool, err error)
	Gt() (result bool, err error)
	Gte() (result bool, err error)
	Lt() (result bool, err error)
	Lte() (result bool, err error)
}

// GCV - Generic Comparable Values
type GCV[T ComparableValueTypes] struct {
	V1       T
	V2       T
	EqError  error
	NeqError error
	GtError  error
	GteError error
	LtError  error
	LteError error
}

func NewGCV[T ComparableValueTypes](v1 T, v2 T) *GCV[T] {
	return &GCV[T]{
		V1: v1,
		V2: v2,
	}
}

func NewOnlyEqGCV[T ComparableValueTypes](v1 T, v2 T) *GCV[T] {
	return &GCV[T]{
		V1:       v1,
		V2:       v2,
		GtError:  ErrCompareOperationIsNotAcceptable,
		GteError: ErrCompareOperationIsNotAcceptable,
		LtError:  ErrCompareOperationIsNotAcceptable,
		LteError: ErrCompareOperationIsNotAcceptable,
	}
}

func (v GCV[T]) Eq() (bool, error) {
	return v.V1 == v.V2, v.EqError
}
func (v GCV[T]) Neq() (bool, error) {
	return v.V1 != v.V2, v.NeqError
}
func (v GCV[T]) Gt() (bool, error) {
	return v.V1 > v.V2, v.GtError
}
func (v GCV[T]) Gte() (bool, error) {
	return v.V1 >= v.V2, v.GteError
}
func (v GCV[T]) Lt() (bool, error) {
	return v.V1 < v.V2, v.LtError
}
func (v GCV[T]) Lte() (bool, error) {
	return v.V1 <= v.V2, v.LteError
}
