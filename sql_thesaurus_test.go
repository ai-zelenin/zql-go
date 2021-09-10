package zql

import (
	"fmt"
)

func ExampleSQLThesaurus_FilterToWherePart() {
	code := `"age">=15 || ("age"<=17 && "name">"fuck" && "name"==nil)`
	q, err := Run(code, NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	sqlt := NewSQLThesaurus("postgres")
	wherePart, args, err := sqlt.FilterToWherePart(q.Filter, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
	fmt.Println(wherePart, args)
}
