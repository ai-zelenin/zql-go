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
	sqlt := NewSQLThesaurus("postgres")
	wherePart, args, err := sqlt.FilterToWherePart(q.Filter, true)
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
	q.With = map[string]*Query{}
	q.With["with_test"] = &Query{
		Model: "test",
	}
	q.Join = append(q.Join, Join{
		Type:      JoinTypeInner,
		Table:     "joined_table",
		Predicate: Eq("joined_table.id", "table.id"),
	})
	sqlt := NewSQLThesaurus("postgres")
	sqlQuery, args, err := sqlt.ToSQL(q, false, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
	fmt.Println(sqlQuery, args)
	expected := `WITH with_test AS (SELECT * FROM "test") SELECT * FROM * INNER JOIN "joined_table" ON ("joined_table"."id" = "table"."id") WHERE (("age" >= 15) OR (("age" <= 17) AND ("name" > 'fuck') AND ("name" IS NULL))) LIMIT 50`
	if sqlQuery != expected {
		t.Fail()
	}
}
