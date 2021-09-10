# ZQL

### Install

```bash
go get -u github.com/ai-zelenin/zql-go
```

### Usage

#### Build Query from expr

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

#### Make SQL from query

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

#### Make only where part (useful for orm)

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
          "field": "f0",
          "op": "gte",
          "value": 0
        },
        {
          "op": "and",
          "value": [
            {
              "field": "f1",
              "op": "eq",
              "value": 1
            },
            {
              "field": "f2",
              "op": "eq",
              "value": 2
            },
            {
              "field": "f3",
              "op": "eq",
              "value": 3
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

#### Use Validator

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

#### Create custom validation

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

#### Create custom predicate Op

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