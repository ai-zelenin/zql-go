package zql

import "fmt"

type ErrCode int

const (
	ErrCodeFieldUnacceptable             = 901
	ErrCodeValueTypeUnacceptableForField = 902
	ErrCodeUnacceptableOp                = 903
	ErrCodeValueTypeUnacceptableForOp    = 904
	ErrCodeValueKindUnacceptableForOp    = 905
	ErrCodeTooManyPredicatesInFilter     = 906
	ErrCodeExceedMaxValueSize            = 907
)

type Error struct {
	Origin error
	Code   ErrCode
}

func NewError(origin error, code ErrCode) *Error {
	if origin == nil {
		return nil
	}
	return &Error{Origin: origin, Code: code}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v code:%d", e.Origin, e.Code)
}
