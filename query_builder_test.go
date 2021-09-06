package zql

import "fmt"

func ExampleQueryBuilder_Build() {
	qb := NewQueryBuilder()
	and := And(
		Lte("age", 17),
		Gt("age", 17),
		Eq("name", nil),
	)
	or := Or(Gte("age", 15), and)
	qb.Filter(or).Page(10, 50)
	q := qb.Build()
	fmt.Println(q)
}
