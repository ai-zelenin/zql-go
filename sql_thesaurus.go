package zql

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"strings"
)

type OpToSQLFunc func(field string, value interface{}) (goqu.Expression, error)

type SQLThesaurus struct {
	dialect       string
	opFunctionMap map[string]OpToSQLFunc
}

func NewSQLThesaurus(dialect string) *SQLThesaurus {
	t := &SQLThesaurus{
		dialect: dialect,
		opFunctionMap: map[string]OpToSQLFunc{
			EQ: func(field string, value interface{}) (goqu.Expression, error) {
				if isNil(value) {
					return goqu.I(field).IsNull(), nil
				}
				return goqu.I(field).Eq(value), nil
			},
			NEQ: func(field string, value interface{}) (goqu.Expression, error) {
				if isNil(value) {
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

func (s *SQLThesaurus) SetOpFunc(key string, f OpToSQLFunc) {
	s.opFunctionMap[key] = f
}

func (s *SQLThesaurus) QueryToSQL(q *Query, prepared bool) (string, []interface{}, error) {
	predicates := q.Filter
	if len(predicates) == 0 {
		return "", nil, nil
	}
	fields := make([]interface{}, len(q.Fields))
	for i, f := range q.Fields {
		fields[i] = f
	}
	if len(fields) == 0 {
		fields = []interface{}{"*"}
	}

	model := "*"
	if q.Model != "" {
		model = q.Model
	}

	dialect := goqu.Dialect(s.dialect)
	query := dialect.Select(fields...).From(model)
	root := goqu.And()
	for _, predicate := range predicates {
		expr, err := s.PredicateToExpression(predicate)
		if err != nil {
			return "", nil, err
		}
		root = root.Append(expr)
	}
	query = query.Where(root)

	if q.Uniq != "" {
		query = query.Distinct(q.Uniq)
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
	return query.Prepared(prepared).ToSQL()
}

func (s *SQLThesaurus) FilterToExpr(predicates []*Predicate) (goqu.Expression, error) {
	if len(predicates) == 0 {
		return nil, nil
	}
	root := goqu.And()
	for _, predicate := range predicates {
		expr, err := s.PredicateToExpression(predicate)
		if err != nil {
			return nil, err
		}
		root = root.Append(expr)
	}
	return root, nil
}

func (s *SQLThesaurus) ExprToWherePart(ex goqu.Expression, prepared bool) (string, []interface{}, error) {
	if ex == nil {
		return "", nil, nil
	}
	dialect := goqu.Dialect(s.dialect)
	query := dialect.Select("*")
	query = query.Where(ex)
	sql, values, err := query.Prepared(prepared).ToSQL()
	if err != nil {
		return "", nil, err
	}
	wherePart := strings.TrimSpace(strings.TrimPrefix(sql, "SELECT * WHERE"))
	return wherePart, values, nil
}

func (s *SQLThesaurus) PredicateToExpression(p *Predicate) (goqu.Expression, error) {
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
			expr, err := s.PredicateToExpression(predicate)
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
			expr, err := s.PredicateToExpression(predicate)
			if err != nil {
				return nil, err
			}
			or = or.Append(expr)
		}
		return or, nil
	default:
		f, ok := s.opFunctionMap[op]
		if !ok {
			return nil, fmt.Errorf("unknown operator %s", op)
		}
		return f(field, value)
	}
}
