package zql

type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)

type Order struct {
	Field     string    `json:"field,omitempty" yaml:"field"`
	Direction Direction `json:"direction,omitempty" yaml:"direction"`
}

func NewOrder(field string, dir Direction) *Order {
	return &Order{
		Field:     field,
		Direction: dir,
	}
}

func Asc(field string) *Order {
	return NewOrder(field, ASC)
}

func Desc(field string) *Order {
	return NewOrder(field, DESC)
}
