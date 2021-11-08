package zql

type JoinType string

const (
	JoinTypeInner JoinType = "inner"
	JoinTypeLeft  JoinType = "left"
	JoinTypeRight JoinType = "right"
)

type Join struct {
	Type      JoinType   `json:"type,omitempty" yaml:"type"`
	Table     string     `json:"table,omitempty" yaml:"table"`
	Predicate *Predicate `json:"predicate,omitempty" yaml:"predicate"`
}
