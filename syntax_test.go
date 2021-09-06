package zql

import (
	"fmt"
)

func ExampleRun() {
	code := `"age">=15 || ("age"<=17 && "name">"fuck" && "name"==nil)`
	q, err := Run(code, NewSyntaxV1())
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
}
