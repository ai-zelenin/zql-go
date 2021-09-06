package zql

import (
	"fmt"
)

func ExampleSQLThesaurus_FilterToSQL() {
	code := `"age">=15 || ("age"<=17 && "name">"fuck" && "name"==nil)`
	q, err := Run(code, NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	sqlt := NewSQLThesaurus("postgres")
	cond, args, err := sqlt.FilterToSQL(q.Filter)
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
	fmt.Println(cond, args)
}
