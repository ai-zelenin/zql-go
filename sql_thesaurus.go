package zql

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"strings"
)

type PredicateOpConvertFunc func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error)
type FieldConvertFunc func(fieldName string) FieldExp

type FieldExp interface {
	exp.Comparable
	exp.Isable
	exp.Inable
	exp.Likeable
	exp.Rangeable
	exp.Orderable
}

type SQLThesaurus struct {
	dialect          string
	opFuncMap        map[string]PredicateOpConvertFunc
	fieldConvertFunc FieldConvertFunc
}

func NewSQLThesaurus(dialect string, fieldConvertFunc FieldConvertFunc) *SQLThesaurus {
	if fieldConvertFunc == nil {
		fieldConvertFunc = func(fieldName string) FieldExp {
			return goqu.I(Sanitize(fieldName))
		}
	}
	t := &SQLThesaurus{
		dialect:          dialect,
		fieldConvertFunc: fieldConvertFunc,
		opFuncMap: map[string]PredicateOpConvertFunc{
			AND: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				predicates, ok := value.([]*Predicate)
				if !ok {
					return nil, fmt.Errorf("%s op must contain predicats slice in value", AND)
				}
				var and = goqu.And()
				for _, predicate := range predicates {
					expr, err := t.PredicateToExpression(predicate, false)
					if err != nil {
						return nil, err
					}
					and = and.Append(expr)
				}
				return and, nil
			},
			OR: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				predicates, ok := value.([]*Predicate)
				if !ok {
					return nil, fmt.Errorf("%s op must contain predicats slice in value", OR)
				}
				var or = goqu.Or()
				for _, predicate := range predicates {
					expr, err := t.PredicateToExpression(predicate, false)
					if err != nil {
						return nil, err
					}
					or = or.Append(expr)
				}
				return or, nil
			},
			EQ: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				if IsNilValue(value) {
					return t.fieldConvertFunc(field).IsNull(), nil
				}
				return t.fieldConvertFunc(field).Eq(value), nil
			},
			NEQ: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				if IsNilValue(value) {
					return t.fieldConvertFunc(field).IsNotNull(), nil
				}
				return t.fieldConvertFunc(field).Neq(value), nil
			},
			GT: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				return t.fieldConvertFunc(field).Gt(value), nil
			},
			GTE: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				return t.fieldConvertFunc(field).Gte(value), nil
			},
			LT: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				return t.fieldConvertFunc(field).Lt(value), nil
			},
			LTE: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				return t.fieldConvertFunc(field).Lte(value), nil
			},
			IN: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				return t.fieldConvertFunc(field).In(value), nil
			},
			LIKE: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				return t.fieldConvertFunc(field).Like(value), nil
			},
			ILIKE: func(t *SQLThesaurus, field string, value interface{}) (goqu.Expression, error) {
				return t.fieldConvertFunc(field).ILike(value), nil
			},
		},
	}
	return t
}

func (s *SQLThesaurus) SetOpFunc(key string, f PredicateOpConvertFunc) {
	s.opFuncMap[key] = f
}

func (s *SQLThesaurus) PredicateToExpression(p *Predicate, valueAsIdent bool) (goqu.Expression, error) {
	var field = p.Field
	var value = p.Value
	if valueAsIdent {
		v, ok := p.Value.(string)
		if ok {
			value = goqu.I(v)
		}
	}
	var op = strings.ToLower(p.Op)
	f, ok := s.opFuncMap[op]
	if !ok {
		return nil, fmt.Errorf("unknown operator %s", op)
	}
	return f(s, field, value)
}

func (s *SQLThesaurus) FilterToExpression(predicates []*Predicate, asIdent bool) (goqu.Expression, error) {
	if len(predicates) == 0 {
		return nil, nil
	}
	root := goqu.And()
	for _, predicate := range predicates {
		expr, err := s.PredicateToExpression(predicate, asIdent)
		if err != nil {
			return nil, err
		}
		root = root.Append(expr)
	}
	return root, nil
}

func (s *SQLThesaurus) FilterExpressionToWherePart(ex goqu.Expression, prepared bool) (string, []interface{}, error) {
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

func (s *SQLThesaurus) FilterToWherePart(predicates []*Predicate, prepared, asIdent bool) (string, []interface{}, error) {
	ex, err := s.FilterToExpression(predicates, asIdent)
	if err != nil {
		panic(err)
	}
	return s.FilterExpressionToWherePart(ex, prepared)
}

func (s *SQLThesaurus) ToGoQu(q *Query, bounded bool) (*goqu.SelectDataset, error) {
	dialect := goqu.Dialect(s.dialect)
	// Fields
	fields := make([]interface{}, len(q.Fields))
	for i, f := range q.Fields {
		fields[i] = f
	}
	if len(fields) == 0 {
		fields = []interface{}{"*"}
	}
	query := dialect.Select(fields...)

	// Table
	model := "*"
	if q.Model != "" {
		model = q.Model
	}
	query = query.From(model)

	// Distinct
	if q.Uniq != "" {
		query = query.Distinct(q.Uniq)
	}

	// Where
	where, err := s.FilterToExpression(q.Filter, false)
	if err != nil {
		return nil, err
	}
	if where != nil {
		query = query.Where(where)
	}

	// Orders
	for _, order := range q.Orders {
		var dir exp.SortDirection
		var nullType exp.NullSortType
		switch Direction(strings.ToLower(string(order.Direction))) {
		case ASC:
			dir = exp.AscDir
			nullType = exp.NullsFirstSortType
		case DESC:
			dir = exp.DescSortDir
			nullType = exp.NullsLastSortType
		}
		expr := exp.NewOrderedExpression(goqu.L(order.Field), dir, nullType)
		query = query.OrderAppend(expr)
	}

	// Limit Offset
	l, o := q.LimitOffset(bounded)
	query = query.Limit(uint(l))
	query = query.Offset(uint(o))
	return query, nil
}

func (s *SQLThesaurus) ToSQL(q *Query, prepared, bounded bool) (string, []interface{}, error) {
	sds, err := s.ToGoQu(q, bounded)
	if err != nil {
		return "", nil, err
	}
	return sds.Prepared(prepared).ToSQL()
}
