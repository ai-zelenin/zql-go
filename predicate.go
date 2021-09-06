package dao

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
	Field string
	Op    string
	Value interface{}
}

func NewPredicate(op string, field string, value interface{}) *Predicate {
	return &Predicate{
		Field: field,
		Value: value,
		Op:    op,
	}
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
