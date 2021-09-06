package dao

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"reflect"
	"strings"
)

type OpToSQLFunc func(field string, value interface{}) (goqu.Expression, error)

type ThesaurusSQL struct {
	dialect       string
	opFunctionMap map[string]OpToSQLFunc
}

func NewThesaurusSQL(dialect string) *ThesaurusSQL {
	t := &ThesaurusSQL{
		dialect: dialect,
		opFunctionMap: map[string]OpToSQLFunc{
			EQ: func(field string, value interface{}) (goqu.Expression, error) {
				if reflect.ValueOf(value).IsNil() {
					return goqu.I(field).IsNull(), nil
				}
				return goqu.I(field).Eq(value), nil
			},
			NEQ: func(field string, value interface{}) (goqu.Expression, error) {
				if reflect.ValueOf(value).IsNil() {
					return goqu.I(field).IsNotNull(), nil
				}
				return goqu.I(field).Neq(value), nil
			},
			GT: func(field string, value interface{}) (goqu.Expression, error) {
				return goqu.I(field).Gt(value), nil
			},
			GTE: func(field string, value interface{}) (goqu.Expression, error) {
				return goqu.I(field).Gte(value), nil
			},
			LT: func(field string, value interface{}) (goqu.Expression, error) {
				return goqu.I(field).Lt(value), nil
			},
			LTE: func(field string, value interface{}) (goqu.Expression, error) {
				return goqu.I(field).Lte(value), nil
			},
			IN: func(field string, value interface{}) (goqu.Expression, error) {
				return goqu.I(field).In(value), nil
			},
			LIKE: func(field string, value interface{}) (goqu.Expression, error) {
				return goqu.I(field).Like(value), nil
			},
			ILIKE: func(field string, value interface{}) (goqu.Expression, error) {
				return goqu.I(field).ILike(value), nil
			},
		},
	}
	return t
}

func (t *ThesaurusSQL) SetOpFunc(key string, f OpToSQLFunc) {
	t.opFunctionMap[key] = f
}

func (t *ThesaurusSQL) QueryToSQL(q *Query) (string, []interface{}, error) {
	predicates := q.Filter
	if len(predicates) == 0 {
		return "", nil, nil
	}
	dialect := goqu.Dialect(t.dialect)
	query := dialect.Select("*")

	root := goqu.And()
	for _, predicate := range predicates {
		expr, err := t.PredicateToExpression(predicate)
		if err != nil {
			return "", nil, err
		}
		root = root.Append(expr)
	}
	query = query.Where(root)

	if q.Distinct != "" {
		query = query.Distinct(q.Distinct)
	}
	for _, order := range q.Orders {
		var dir exp.SortDirection
		var nullType exp.NullSortType
		switch order.Direction {
		case ASC:
			dir = exp.AscDir
			nullType = exp.NullsFirstSortType
		case DESC:
			dir = exp.DescSortDir
			nullType = exp.NullsLastSortType
		}
		expr := exp.NewOrderedExpression(goqu.I(order.Field), dir, nullType)
		query = query.OrderAppend(expr)
	}

	l, o := q.LimitOffset()
	query = query.Limit(uint(l))
	query = query.Offset(uint(o))
	return query.Prepared(true).ToSQL()
}

func (t *ThesaurusSQL) FilterToSQL(predicates []*Predicate) (string, []interface{}, error) {
	if len(predicates) == 0 {
		return "", nil, nil
	}
	dialect := goqu.Dialect(t.dialect)
	root := goqu.And()
	for _, predicate := range predicates {
		expr, err := t.PredicateToExpression(predicate)
		if err != nil {
			return "", nil, err
		}
		root = root.Append(expr)
	}
	query := dialect.Select("*")
	query.Where(root)
	sql, values, err := query.Prepared(true).ToSQL()
	if err != nil {
		return "", nil, err
	}
	wherePart := strings.TrimSpace(strings.TrimPrefix(sql, "SELECT * WHERE"))
	return wherePart, values, nil
}

func (t *ThesaurusSQL) PredicateToExpression(p *Predicate) (goqu.Expression, error) {
	var field = p.Field
	var value = p.Value
	var op = strings.ToLower(p.Op)
	switch op {
	case AND:
		predicates, ok := value.([]*Predicate)
		if !ok {
			return nil, fmt.Errorf("%s op must contain predicats slice in value", op)
		}
		var and = goqu.And()
		for _, predicate := range predicates {
			expr, err := t.PredicateToExpression(predicate)
			if err != nil {
				return nil, err
			}
			and = and.Append(expr)
		}
		return and, nil
	case OR:
		predicates, ok := value.([]*Predicate)
		if !ok {
			return nil, fmt.Errorf("%s op must contain predicats slice in value", op)
		}
		var or = goqu.Or()
		for _, predicate := range predicates {
			expr, err := t.PredicateToExpression(predicate)
			if err != nil {
				return nil, err
			}
			or = or.Append(expr)
		}
		return or, nil
	default:
		f, ok := t.opFunctionMap[op]
		if !ok {
			return nil, fmt.Errorf("unknown operator %s", op)
		}
		return f(field, value)
	}
}
