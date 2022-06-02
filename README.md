# ZQL

### Install

```bash
go get -u github.com/ai-zelenin/zql-go
```

### Concept

ZQL is machine friendly abstract query format.
It used for machine processing/generating data queries.

Almost any data query can be represented as predicate tree.  
A Predicate in programming is an expression that uses one or more values with a boolean result.  
ZQL Predicate looks like this:

##### Golang

```go
type Predicate struct {
    Field string   
    Op    string    
    Value interface{}
}
```

##### Json example

```json
{
  "field": "field_name",
  "op": "gte",
  "value": 18
}
```

For example, we have table "humans"
<table>
<thead>
<tr>
<th>id</th>
<th>name</th>
<th>age</th>
<th>sex</th>
<th>status</th>
<th>created_at</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td>pole</td>
<td>16</td>
<td>1</td>
<td>1</td>
<td>2022-06-01T01:01:30</td>
</tr>
<tr>
<td>2</td>
<td>nickol</td>
<td>42</td>
<td>2</td>
<td>1</td>
<td>2022-06-02T02:03:30</td>
</tr>
<tr>
<td>3</td>
<td>ivan</td>
<td>34</td>
<td>1</td>
<td>2</td>
<td>2022-05-30T08:06:10</td>
</tr>
</tbody>
</table>

We want select only humans which is male and adult or female and have status=2

Such SQL query look like this:

```sql
select *
from "humans"
where ("sex" = 1 AND "age" >= 18)
   OR ("sex" = 2 AND "status" = 2); 
```

This query contains 7 predicates.  
4 with arithmetic operations.  
3 with logic operations.  

```text
P1: sex == 1
P2: age >= 18
P3: P1 && P2
P4: sex == 2
P5: status == 2
P6: P4 && P5
P7: P3 || P6
```

If translate SQL to ZQL it looks like this
```json
{
  "filter": [
    {
      "op": "or",
      "value": [
        {
          "op": "and",
          "value": [
            {
              "field": "sex",
              "op": "eq",
              "value": 1
            },
            {
              "field": "age",
              "op": "gte",
              "value": 18
            }
          ]
        },
        {
          "op": "and",
          "value": [
            {
              "field": "sex",
              "op": "eq",
              "value": 2
            },
            {
              "field": "status",
              "op": "eq",
              "value": 2
            }
          ]
        }
      ]
    }
  ]
}
```

By default zql-go can convert ZQL data queries to SQL data queries 
with this predicate operations

#### Logic
* AND ("and")
* OR ("or")


#### Arithmetic
* EQ ("eq")
* GT ("gt")
* GTE ("gte")
* LT ("lt")
* LTE ("lte")
* NEQ ("neq")
* IN ("in")
* LIKE ("like")
* ILIKE ("ilike")

#### Text
* LIKE ("like")    
* ILIKE ("ilike")   


#### Examples

##### Use Query builder

```go
package main

import (
	"fmt"
	"github.com/ai-zelenin/zql-go"
)

func main() {
	qb := zql.NewQueryBuilder()
	qb.Filter(
		zql.Or(
			zql.Gte("f0", 0),
			zql.And(
				zql.Eq("f1", 1),
				zql.Eq("f2", 2),
				zql.Eq("f3", 3),
			)),
	)
	qb.Page(10, 5)
	q := qb.Build()

	sqlt := zql.NewSQLThesaurus("postgres")
	sql, args, err := sqlt.QueryToSQL(q, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(q)
	fmt.Println(sql)
	fmt.Println(args)
}
```

##### Use expr

```go

package main

import (
	"fmt"
	"github.com/ai-zelenin/zql-go"
)

func main() {
	code := `"age">=18 || ("age"<18 && "name"=="lol")`
	q, err := zql.Run(code, zql.NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
}
```

##### Build only where part (useful for orm)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/ai-zelenin/zql-go"
)

func main() {
	var queryString = `
{
  "filter": [
    {
      "op": "or",
      "value": [
        {
          "op": "and",
          "value": [
            {
              "field": "sex",
              "op": "eq",
              "value": 1
            },
            {
              "field": "age",
              "op": "gte",
              "value": 18
            }
          ]
        },
        {
          "op": "and",
          "value": [
            {
              "field": "sex",
              "op": "eq",
              "value": 2
            },
            {
              "field": "status",
              "op": "eq",
              "value": 2
            }
          ]
        }
      ]
    }
  ]
}
`
	q := new(zql.Query)
	err := json.Unmarshal([]byte(queryString), q)
	if err != nil {
		panic(err)
	}

	sqlt := zql.NewSQLThesaurus("postgres")
	wherePart, args, err := sqlt.FilterToWherePart(q.Filter, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(q)
	fmt.Println(wherePart)
	fmt.Println(args)
}
```

##### Use Validator

```go

package main

import (
	"github.com/ai-zelenin/zql-go"
	"time"
)

type Status int

type Human struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Age       int64     `json:"age"`
	Sex       bool      `json:"sex"`
	BirthDate time.Time `json:"birth_date"`
}

func main() {
	code := `"age">="18" || ("age"<18 && "name"=="lol")`
	q, err := zql.Run(code, zql.NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	var validator = zql.NewExtendableValidator()
	validator.SetupValidatorForModel(new(Human), "json")
	err = validator.Validate(q)
	if err != nil {
		panic(err)
	}
}
```

##### Create custom validation

```go

package main

import (
	"fmt"
	"github.com/ai-zelenin/zql-go"
	"reflect"
)

func main() {
	code := `"age">=18 || ("age"<18 && "name"=="admin")`
	q, err := zql.Run(code, zql.NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	var validator = zql.NewExtendableValidator()
	var customValidation zql.ValidatePredicateFunc = func(field, op string, value interface{}, rt reflect.Type, rv reflect.Value) error {
		if field == "name" && rv.String() == "admin" {
			return zql.NewError(fmt.Errorf("forbidden value for field %s", field), zql.ErrCode(1901))
		}
		return nil
	}
	validator.AddPredicateValidator(customValidation)
	err = validator.Validate(q)
	if err != nil {
		panic(err)
	}
}
```

##### Create custom predicate Op

```go

package main

import (
	"fmt"
	"github.com/ai-zelenin/zql-go"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"reflect"
)

func main() {
	qb := zql.NewQueryBuilder()
	qb.Filter(
		zql.Or(
			zql.Gte("f0", 0),
			zql.And(
				zql.Eq("f1", 1),
				zql.Eq("f2", 2),
				zql.Eq("f3", 3),
			)),
		zql.NewPredicate("between", "id", []int{1, 3}),
	)
	qb.Page(10, 5)
	q := qb.Build()

	sqlt := zql.NewSQLThesaurus("postgres")

	// New Op
	sqlt.SetOpFunc("between", func(t *zql.SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			from := rv.Index(0).Interface()
			to := rv.Index(1).Interface()
			return goqu.I(field).Between(exp.NewRangeVal(from, to)), nil
		default:
			return nil, fmt.Errorf("bad value for op between")
		}
	})

	sql, args, err := sqlt.QueryToSQL(q, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(q)
	fmt.Println(sql)
	fmt.Println(args)
}
```