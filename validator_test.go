package zql

import (
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"reflect"
	"testing"
	"time"
)

type Status int

type MetaData struct {
	Field1 string
	Field2 int
	Field3 bool
}

type Human struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Age       int64     `json:"age"`
	Sex       bool      `json:"sex"`
	BirthDate time.Time `json:"birth_date"`
	Meta      *MetaData `json:"meta"`
}

type ValidatorTestCase struct {
	Code      string
	ErrorCode ErrCode
}

var ValidatorTestTable = []ValidatorTestCase{
	{
		Code: `
("sex" == true && "age" >= 18) ||
("sex" == false && "age" >= 21) ||
"birth_date" < now() ||
"id" in [1,2,3] ||
"age" == nil ||
"status" != 1 || 
like("name","%Lu%") ||
ilike("name","%lu%")
`,
		ErrorCode: 0,
	},
	{
		Code:      `("InvalidFieldName" == true && "age" >= 18)`,
		ErrorCode: ErrCodeFieldUnacceptable,
	},
	{
		Code:      `("sex" == "InvalidCompareValueType" && "age" >= 18)`,
		ErrorCode: ErrCodeValueTypeUnacceptableForField,
	},
	{
		Code:      `qb.Filter(("sex" == true && "age" >= 18)).Order(asc("id"))`,
		ErrorCode: 0,
	},
}

func TestExtendableValidator_Validate(t *testing.T) {
	for _, testCase := range ValidatorTestTable {
		q, err := Run(testCase.Code, NewSyntaxV1())
		if err != nil {
			panic(err)
		}
		code := ErrCode(0)
		err = validate(q)
		if err != nil {
			e, ok := err.(*Error)
			if ok {
				code = e.Code
			} else {
				code = ErrCode(-1)
			}
		}
		if code != testCase.ErrorCode {
			t.Fatalf("'%s' -> bad err code -  have:%v expect:%v", testCase.Code, code, testCase.ErrorCode)
		}
	}
}

func validate(q *Query) error {
	var validator = NewExtendableValidator()
	model := make([]*Human, 0)
	validator.SetupValidatorForModel(model, "json")
	sqlt := NewSQLThesaurus("postgres")
	sqlt.SetOpFunc("between", func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			from := rv.Index(0).Interface()
			to := rv.Index(1).Interface()
			return goqu.I(field).Between(exp.NewRangeVal(from, to)), nil
		default:
			return nil, fmt.Errorf("bad value for op between")
		}
	})
	err := validator.Validate(q)
	if err != nil {
		return err
	}

	sql, _, err := sqlt.ToSQL(q, false, true)
	if err != nil {
		return err
	}

	data, err := json.Marshal(q)
	if err != nil {
		return err
	}

	afterSerializeQuery := new(Query)
	err = json.Unmarshal(data, afterSerializeQuery)
	if err != nil {
		return err
	}

	err = validator.Validate(afterSerializeQuery)
	if err != nil {
		return err
	}
	sqlAS, _, err := sqlt.ToSQL(afterSerializeQuery, false, true)
	if err != nil {
		return err
	}
	if sql != sqlAS {
		fmt.Println(sql)
		fmt.Println(sqlAS)
		return NewError(fmt.Errorf("serialze side effect"), ErrCode(-2))
	}
	return nil
}
