package zql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleRun() {
	code := `
	("user.age">=18) || 
	(
		"user.age"<18 &&
		"user.sex"=="male" &&
		"user.location"!=nil
	) || 
	("location.lat" > 57.1 && "location.lon" < 39.2)`
	q, err := Run(code, NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
	/*
		{
			"filter": [
				{
					"op": "or",
					"value": [
						{
							"field": "user.age",
							"op": "gte",
							"value": 18
						},
						{
							"op": "and",
							"value": [
								{
									"field": "user.age",
									"op": "lt",
									"value": 18
								},
								{
									"field": "user.sex",
									"op": "eq",
									"value": "male"
								},
								{
									"field": "user.location",
									"op": "neq",
									"value": null
								}
							]
						},
						{
							"op": "and",
							"value": [
								{
									"field": "location.lat",
									"op": "gt",
									"value": 57.1
								},
								{
									"field": "location.lon",
									"op": "lt",
									"value": 39.2
								}
							]
						}
					]
				}
			]
		}
	*/
}

func TestSyntaxRun(t *testing.T) {
	code := `
	("user.age">=18) || 
	(
		"user.age"<18 &&
		"user.sex"=="male" &&
		"user.location"!=nil
	) || 
	("location.lat" > 57.1 && "location.lon" < 39.2)`
	q, err := Run(code, NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	expected := `
{
	"filter": [
		{
			"op": "or",
			"value": [
				{
					"field": "user.age",
					"op": "gte",
					"value": 18
				},
				{
					"op": "and",
					"value": [
						{
							"field": "user.age",
							"op": "lt",
							"value": 18
						},
						{
							"field": "user.sex",
							"op": "eq",
							"value": "male"
						},
						{
							"field": "user.location",
							"op": "neq",
							"value": null
						}
					]
				},
				{
					"op": "and",
					"value": [
						{
							"field": "location.lat",
							"op": "gt",
							"value": 57.1
						},
						{
							"field": "location.lon",
							"op": "lt",
							"value": 39.2
						}
					]
				}
			]
		}
	]
}`
	fmt.Println(q)
	assert.JSONEq(t, expected, q.String())
}
