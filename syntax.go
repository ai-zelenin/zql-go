package zql

import (
	"fmt"
	"github.com/antonmedv/expr"
)

const MaxPerPage = 1000
const DefaultPerPage = 50

type Operator struct {
	OperatorName string
	Func         interface{}
}

func NewOperator(operatorName string, f interface{}) *Operator {
	return &Operator{OperatorName: operatorName, Func: f}
}

type Syntax interface {
	Env() map[string]interface{}
}

func Run(code string, syntax Syntax) (*Query, error) {
	env := syntax.Env()
	var options = []expr.Option{expr.Env(env)}
	for k, v := range env {
		o, ok := v.(*Operator)
		if ok {
			options = append(options, expr.Operator(o.OperatorName, k))
			env[k] = o.Func
		}
	}
	script, err := expr.Compile(code, options...)
	if err != nil {
		return nil, err
	}
	result, err := expr.Run(script, env)
	if err != nil {
		return nil, err
	}
	switch v := result.(type) {
	case *Predicate:
		return &Query{
			Filter: []*Predicate{v},
		}, nil
	case *Query:
		return v, nil
	case *QueryBuilder:
		return v.Build(), nil
	default:
		return nil, fmt.Errorf("unexpected result [type:%T value:%v] from script", result, result)
	}
}
