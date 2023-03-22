package zql

type Direction string

const (
	ASC  Direction = "asc"
	DESC Direction = "desc"
)

type NullSortType string

const (
	NullsFirst NullSortType = "first"
	NullsLast  NullSortType = "last"
)

type OrderOption func(o *Order)

func WithNullsLast() OrderOption {
	return func(o *Order) {
		o.NullSortType = NullsLast
	}
}

type Order struct {
	Field        string       `json:"field,omitempty" yaml:"field"`
	Direction    Direction    `json:"direction,omitempty" yaml:"direction"`
	NullSortType NullSortType `json:"null_sort_type" yaml:"null_sort_type"`
}

func NewOrder(field string, dir Direction, options ...OrderOption) *Order {
	o := &Order{
		Field:     field,
		Direction: dir,
	}
	for _, option := range options {
		option(o)
	}
	return o
}

func Asc(field string, options ...OrderOption) *Order {
	return NewOrder(field, ASC, options...)
}

func Desc(field string, options ...OrderOption) *Order {
	return NewOrder(field, DESC, options...)
}
