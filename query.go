package zql

type Query struct {
	Filter    []*Predicate `json:"filter" yaml:"filter"`
	Orders    []*Order     `json:"orders" yaml:"orders"`
	Relations []string     `json:"relations" yaml:"relations"`
	Page      int64        `json:"page" yaml:"page"`
	PerPage   int64        `json:"per_page" yaml:"per_page"`
	Distinct  string       `json:"distinct" yaml:"distinct"`
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
