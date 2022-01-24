package zql

import (
	"fmt"
	"reflect"
)

type ValidatorPredicateValues struct {
	RequiredFieldValueTypeMap map[string]string
	MaxValueSize              int
}

func NewValidatorPredicateValues() *ValidatorPredicateValues {
	return &ValidatorPredicateValues{
		RequiredFieldValueTypeMap: make(map[string]string),
		MaxValueSize:              -1,
	}
}

func (v *ValidatorPredicateValues) AddFieldValuePair(f FieldDesc) {
	v.RequiredFieldValueTypeMap[f.Name] = f.Type
}

func (v *ValidatorPredicateValues) Validate(field, op string, _ interface{}, rt reflect.Type, rv reflect.Value) error {
	if rt == nil {
		return nil
	}
	valueType := ValueTypeToString(rt)
	requiredValueType, ok := v.RequiredFieldValueTypeMap[field]
	if ok && IsCompareOp(op) {
		isValueTypeGoodForField := requiredValueType == valueType
		if !isValueTypeGoodForField {
			return NewError(fmt.Errorf("value type %v unacceptable for field %s[%s]", valueType, field, requiredValueType), ErrCodeValueTypeUnacceptableForField)
		}
	}

	if v.MaxValueSize > 0 {
		switch rt.Kind() {
		case reflect.Slice, reflect.Array, reflect.String:
			if rv.Len() > v.MaxValueSize {
				return NewError(fmt.Errorf("exceed max value size"), ErrCodeExceedMaxValueSize)
			}
		}
	}
	return nil
}
