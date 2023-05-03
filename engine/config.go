package engine

import (
	"github.com/ai-zelenin/zql-go/compare"
)

type Config struct {
	Comparator         Comparator
	PreHooks           []HookFunc
	PostHooks          []HookFunc
	CompileFunc        CompileFunc
	PredicateFunctions map[string]ExecuteFunc
}

var DefaultConfig = &Config{
	Comparator:  compare.NewTypedComparator(),
	CompileFunc: CompileQLPredicateToExecutable,
}

type CompileFunc func(funcs map[string]ExecuteFunc, p Predicate, preHooks, postHooks []HookFunc) (ep *ExecutablePredicate, err error)

type Comparator interface {
	Eq(v1, v2 any) (result bool, err error)
	Neq(v1, v2 any) (result bool, err error)
	Gt(v1, v2 any) (result bool, err error)
	Gte(v1, v2 any) (result bool, err error)
	Lt(v1, v2 any) (result bool, err error)
	Lte(v1, v2 any) (result bool, err error)
}

type HookFunc func(ep *ExecutablePredicate, scope Scope) error
