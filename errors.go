package zql

import "fmt"

type ErrCode int

const ErrCodeFieldUnacceptable = 901
const ErrCodeValueTypeUnacceptableForField = 902
const ErrCodeValueTypeUnacceptableForOp = 903

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
