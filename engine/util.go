package engine

import (
	"reflect"
	"strings"
)

func CompileQLPredicateToExecutable(funcs map[string]ExecuteFunc, p Predicate, preHooks, postHooks []HookFunc) (ep *ExecutablePredicate, err error) {
	executeFunc, ok := funcs[p.FuncName()]
	if !ok {
		return nil, ErrOpNotFound
	}
	args := p.Args()
	field, ok := args["field"].(string)
	if !ok {
		return nil, ErrBadValueType
	}
	v := args["value"]
	if p.FuncName() == "and" || p.FuncName() == "or" {
		rv := reflect.ValueOf(v)
		n := rv.Len()
		epl := make([]*ExecutablePredicate, n)
		for i := 0; i < n; i++ {
			predicate, isPred := rv.Index(i).Interface().(Predicate)
			if !isPred {
				return nil, ErrBadValueType
			}
			epl[i], err = CompileQLPredicateToExecutable(funcs, predicate, preHooks, postHooks)
			if err != nil {
				return nil, err
			}
		}
		return &ExecutablePredicate{
			FuncName: p.FuncName(),
			Func:     executeFunc,
			Value:    epl,
			PreHook:  preHooks,
			PostHook: postHooks,
		}, nil
	}
	return &ExecutablePredicate{
		FuncName: p.FuncName(),
		Func:     executeFunc,
		Value:    v,
		Field:    field,
		PreHook:  preHooks,
		PostHook: postHooks,
	}, nil
}

func NewVarPrefixHook(varPrefix string) HookFunc {
	return func(ep *ExecutablePredicate, scope Scope) error {
		valueStr, ok := ep.Value.(string)
		if ok {
			if strings.HasPrefix(valueStr, varPrefix) {
				key := strings.TrimPrefix(valueStr, varPrefix)
				val, ok := scope.Get(key)
				if !ok {
					return NewErrVariableNotFound(valueStr)
				}
				ep.Value = val
			}
		}
		return nil
	}
}
