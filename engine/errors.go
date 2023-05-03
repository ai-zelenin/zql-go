package engine

import (
	"errors"
	"fmt"
)

var (
	ErrBadValueType     = errors.New("bad value type")
	ErrOpNotFound       = errors.New("operation not found")
	ErrVariableNotFound = errors.New("variable not found in scope")
)

func NewErrVariableNotFound(varName string) error {
	return fmt.Errorf("%s %w", varName, ErrVariableNotFound)
}
