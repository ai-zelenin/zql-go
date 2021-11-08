package zql

import (
	"fmt"
	"testing"
)

type MutatorTestCase struct {
	Code      string
	ErrorCode ErrCode
}

var MutatorTestTable = []MutatorTestCase{
	{
		Code:      `("InvalidFieldName" == true && "age" >= 18)`,
		ErrorCode: 0,
	},
}

func TestExtendableMutator_Mutate(t *testing.T) {
	for _, testCase := range MutatorTestTable {
		q, err := Run(testCase.Code, NewSyntaxV1())
		if err != nil {
			panic(err)
		}
		code := ErrCode(0)
		err = mutate(q)
		if err != nil {
			e, ok := err.(*Error)
			if ok {
				code = e.Code
			} else {
				code = ErrCode(-1)
			}
		}
		if code != testCase.ErrorCode {
			t.Fatalf("bad err code -  have:%v expect:%v", code, testCase.ErrorCode)
		}
	}
}

func mutate(q *Query) error {
	var mutator = NewExtendableMutator()
	mutator.AddPredicateMutator(MutatePredicateFunc(func(p *Predicate) error {
		if p.Field == "InvalidFieldName" {
			p.Value = 1
		}
		return nil
	}))
	err := mutator.Mutate(q)
	if err != nil {
		panic(err)
	}

	sqlt := NewSQLThesaurus("postgres")
	sql, _, err := sqlt.ToSQL(q, false, true)
	if err != nil {
		return err
	}
	fmt.Println(sql)
	return nil
}
