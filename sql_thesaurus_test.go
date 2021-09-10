package zql

import (
	"fmt"
)

func ExampleSQLThesaurus_ExprToWherePart() {
	code := `"age">=15 || ("age"<=17 && "name">"fuck" && "name"==nil)`
	q, err := Run(code, NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	sqlt := NewSQLThesaurus("postgres")
	ex, err := sqlt.FilterToExpression(q.Filter)
	if err != nil {
		panic(err)
	}
	wherePart, args, err := sqlt.ExprToWherePart(ex, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
	fmt.Println(wherePart, args)
}
