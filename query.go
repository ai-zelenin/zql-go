package zql

type Query struct {
	Model     string       `json:"model,omitempty" yaml:"model"`
	Fields    []string     `json:"fields,omitempty" yaml:"fields"`
	Uniq      string       `json:"uniq,omitempty" yaml:"uniq"`
	Relations []string     `json:"relations,omitempty" yaml:"relations"`
	Filter    []*Predicate `json:"filter,omitempty" yaml:"filter"`
	Orders    []*Order     `json:"orders,omitempty" yaml:"orders"`
	Page      int64        `json:"page,omitempty" yaml:"page"`
	PerPage   int64        `json:"per_page,omitempty" yaml:"per_page"`
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) LimitOffset() (limit int, offset int) {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PerPage == 0 {
		q.PerPage = DefaultPerPage
	}
	if q.PerPage > MaxPerPage {
		q.PerPage = MaxPerPage
	}

	limit = int(q.PerPage)
	offset = int(q.PerPage * (q.Page - 1))
	return
}
