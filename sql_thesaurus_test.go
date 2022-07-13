package zql

import (
	"fmt"
	"testing"
)

func ExampleSQLThesaurus_FilterToWherePart() {
	code := `"age">=15 || ("age"<=17 && "name">"fuck" && "name"==nil)`
	q, err := Run(code, NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	sqlt := NewSQLThesaurus("postgres", nil)
	wherePart, args, err := sqlt.FilterToWherePart(q.Filter, true, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
	fmt.Println(wherePart, args)
}

func TestSQLThesaurus_ToSQL(t *testing.T) {
	code := `"age">=15 || ("age"<=17 && "name">"fuck" && "name"==nil)`
	q, err := Run(code, NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	sqlt := NewSQLThesaurus("postgres", nil)
	sqlQuery, args, err := sqlt.ToSQL(q, false, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
	fmt.Println(sqlQuery, args)
	expected := `SELECT * FROM * WHERE (("age" >= 15) OR (("age" <= 17) AND ("name" > 'fuck') AND ("name" IS NULL))) LIMIT 50`
	if sqlQuery != expected {
		t.Fail()
	}
}
