package compare

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrTypesIsNotComparable            = errors.New("types is not comparable")
	ErrCompareOperationIsNotAcceptable = errors.New("compare operation is not acceptable ")
)

func NewCompareError(rv1, rv2 reflect.Value) error {
	return fmt.Errorf("V1:%v[%s]  V2:%v[%s]  err: %w", rv1, rv1.Kind(), rv2, rv2.Kind(), ErrTypesIsNotComparable)
}
