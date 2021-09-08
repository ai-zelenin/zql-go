package zql

type Query struct {
	Model     string       `json:"model" yaml:"model"`
	Fields    []string     `json:"fields" yaml:"fields"`
	Uniq      string       `json:"uniq" yaml:"uniq"`
	Relations []string     `json:"relations" yaml:"relations"`
	Filter    []*Predicate `json:"filter" yaml:"filter"`
	Orders    []*Order     `json:"orders" yaml:"orders"`
	Page      int64        `json:"page" yaml:"page"`
	PerPage   int64        `json:"per_page" yaml:"per_page"`
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
