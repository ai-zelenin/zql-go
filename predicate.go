package zql

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

const (
	AND   = "and"
	OR    = "or"
	EQ    = "eq"
	GT    = "gt"
	GTE   = "gte"
	LT    = "lt"
	LTE   = "lte"
	NEQ   = "neq"
	IN    = "in"
	LIKE  = "like"
	ILIKE = "ilike"
)

type Predicate struct {
	Field string `json:"field,omitempty" yaml:"field"`
	Op    string `json:"op" yaml:"op"`
	Value any    `json:"value" yaml:"value"`
}

func (p *Predicate) FuncName() string {
	return p.Op
}

func (p *Predicate) Args() map[string]any {
	return map[string]any{
		"op":    p.Op,
		"field": p.Field,
		"value": p.Value,
	}
}

func NewPredicate(op string, field string, value interface{}) *Predicate {
	return &Predicate{
		Field: field,
		Value: value,
		Op:    op,
	}
}

func (p *Predicate) IsGroup() bool {
	return p.Op == AND || p.Op == OR
}

func (p *Predicate) Append(n Node) {
	if p.IsGroup() {
		predicates, ok := p.Value.([]*Predicate)
		if ok {
			predicate, ok := n.(*Predicate)
			if ok {
				p.Value = append(predicates, predicate)
			}
		}
	}
}

func (p *Predicate) Children() []Node {
	var list []Node
	if p.IsGroup() {
		predicates, ok := p.Value.([]*Predicate)
		if ok {
			list = make([]Node, len(predicates))
			for i, predicate := range predicates {
				list[i] = predicate
			}
		}
	}
	return list
}

func (p *Predicate) Walk(cb WalkFunc, parent Node, lvl int) error {
	err := cb(parent, p, lvl)
	if err != nil {
		return err
	}
	for _, node := range p.Children() {
		err = node.Walk(cb, p, lvl+1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Predicate) String() string {
	switch p.Op {
	case AND, OR:
		return fmt.Sprintf("%s(\n%v\n)", p.Op, p.Value)
	default:
		return fmt.Sprintf("%s(%s,%v)", p.Op, p.Field, p.Value)
	}
}

func (p *Predicate) UnmarshalJSON(bytes []byte) error {
	p.Op = gjson.GetBytes(bytes, "op").String()
	if p.Op == AND || p.Op == OR {
		type groupPredicate struct {
			Field string       `json:"field" yaml:"field"`
			Op    string       `json:"op" yaml:"op"`
			Value []*Predicate `json:"value" yaml:"value"`
		}
		pr := new(groupPredicate)
		err := json.Unmarshal(bytes, pr)
		if err != nil {
			return err
		}
		p.Field = pr.Field
		p.Op = pr.Op
		p.Value = pr.Value
	} else {
		type simplePredicate struct {
			Field string      `json:"field" yaml:"field"`
			Op    string      `json:"op" yaml:"op"`
			Value interface{} `json:"value" yaml:"value"`
		}
		pr := new(simplePredicate)
		err := json.Unmarshal(bytes, pr)
		if err != nil {
			return err
		}
		p.Field = pr.Field
		p.Op = pr.Op
		p.Value = pr.Value
	}
	return nil
}

func And(p ...*Predicate) *Predicate {
	return NewPredicate(AND, "", p)
}

func Or(p ...*Predicate) *Predicate {
	return NewPredicate(OR, "", p)
}

func Eq(field string, value interface{}) *Predicate {
	return NewPredicate(EQ, field, value)
}

func Neq(field string, value interface{}) *Predicate {
	return NewPredicate(NEQ, field, value)
}

func Gt(field string, value interface{}) *Predicate {
	return NewPredicate(GT, field, value)
}

func Gte(field string, value interface{}) *Predicate {
	return NewPredicate(GTE, field, value)
}

func Lt(field string, value interface{}) *Predicate {
	return NewPredicate(LT, field, value)
}

func Lte(field string, value interface{}) *Predicate {
	return NewPredicate(LTE, field, value)
}

func In(field string, value interface{}) *Predicate {
	return NewPredicate(IN, field, value)
}

func Like(field string, value interface{}) *Predicate {
	return NewPredicate(LIKE, field, value)
}

func ILike(field string, value interface{}) *Predicate {
	return NewPredicate(ILIKE, field, value)
}
