package zql

import "fmt"

func ExampleQueryBuilder_Build() {
	qb := NewQueryBuilder()
	qb.Filter(
		Or(
			Gte("f0", 0),
			And(
				Eq("f1", 1),
				Eq("f2", 2),
				Eq("f3", 3),
			)),
	)
	qb.Page(10, 5)
	q := qb.Build()
	fmt.Println(q)
}
