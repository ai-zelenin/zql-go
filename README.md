# ZQL

Build Query

```go
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
```

Build Query from expr

```go
code := `"age">=15 || ("age"<=17 && "name">"fuck" && "name"==nil)`
q, err := Run(code, NewSyntaxV1())
if err != nil {
    panic(err)
}
fmt.Println(q)
```

Make SQL where part from Query
```go
code := `"age">=15 || ("age"<=17 && "name">"fuck" && "name"==nil)`
q, err := Run(code, NewSyntaxV1())
if err != nil {
    panic(err)
}
fmt.Println(q)

sqlt := NewSQLThesaurus("postgres")
cond, args, err := sqlt.FilterToSQL(q.Filter)
if err != nil {
    panic(err)
}
fmt.Println(code)
fmt.Println(cond, args)
```