package engine

import (
	"bytes"
	"fmt"
)

type Predicate interface {
	FuncName() string
	Args() map[string]any
}

type ExecuteFunc func(e *ExecutablePredicate, s Scope) (result bool, err error)

type ExecutablePredicate struct {
	Field      string      `json:"field" yaml:"field"`
	Value      any         `json:"value" yaml:"value"`
	Result     bool        `json:"result" yaml:"result"`
	Error      error       `json:"error" yaml:"error"`
	IsExecuted bool        `json:"is_executed" yaml:"is_executed"`
	FuncName   string      `json:"func_name" yaml:"func_name"`
	Func       ExecuteFunc `json:"-" yaml:"-"`
	PreHook    []HookFunc  `json:"-" yaml:"-"`
	PostHook   []HookFunc  `json:"-" yaml:"-"`
}

func (e *ExecutablePredicate) Exec(s Scope) (result bool, err error) {
	e.Reset()
	for _, hookFunc := range e.PreHook {
		if hookFunc != nil {
			err = hookFunc(e, s)
			if err != nil {
				return false, err
			}
		}
	}
	e.Result, e.Error = e.Func(e, s)
	e.IsExecuted = true
	for _, hookFunc := range e.PostHook {
		if hookFunc != nil {
			err = hookFunc(e, s)
			if err != nil {
				return false, err
			}
		}
	}
	return e.Result, e.Error
}

func (e *ExecutablePredicate) Reset() {
	e.Result = false
	e.Error = nil
	e.IsExecuted = false
	v, ok := e.Value.([]*ExecutablePredicate)
	if ok {
		for _, predicate := range v {
			predicate.Reset()
		}
	}
}

func (e *ExecutablePredicate) Dump(prefix string) string {
	b := bytes.NewBuffer(nil)
	b.WriteString(prefix + e.FuncName)
	if e.Field != "" {
		b.WriteString(fmt.Sprintf("(%s,%v)", e.Field, e.Value))
	}
	if e.IsExecuted {
		b.WriteString(fmt.Sprintf(":%t\n", e.Result))
	} else {
		b.WriteString(":skip\n")
	}
	v, ok := e.Value.([]*ExecutablePredicate)
	if ok {
		for _, predicate := range v {
			b.WriteString(predicate.Dump(prefix + "\t"))
		}
	}
	return b.String()
}
