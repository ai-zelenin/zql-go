package zql

import (
	"fmt"
	"testing"
)

type MutatorTestCase struct {
	Query     *Query
	ErrorCode ErrCode
}

var MutatorTestTable = []MutatorTestCase{
	{
		Query: &Query{
			Filter: []*Predicate{
				{
					Op: AND,
					Value: []*Predicate{
						{
							Field: "InvalidFieldName",
							Op:    EQ,
							Value: true,
						},
						{
							Field: "age",
							Op:    GTE,
							Value: 18,
						},
					},
				},
			},
			Relations: map[string]*Query{
				"Rel": {
					Filter: []*Predicate{
						{
							Op: AND,
							Value: []*Predicate{
								{
									Field: "InvalidFieldName",
									Op:    EQ,
									Value: false,
								},
								{
									Field: "age",
									Op:    GTE,
									Value: 18,
								},
							},
						},
					},
				},
			},
		},
		ErrorCode: 0,
	},
}

func TestExtendableMutator_Mutate(t *testing.T) {
	for _, testCase := range MutatorTestTable {
		code := ErrCode(0)
		err := mutate(testCase.Query)
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

	sqlt := NewSQLThesaurus("postgres", nil)
	sql, _, err := sqlt.ToSQL(q, false, true)
	if err != nil {
		return err
	}
	fmt.Println(sql)
	return nil
}
