package zql

import "encoding/json"

type Query struct {
	Model     string            `json:"model,omitempty" yaml:"model"`
	Fields    []string          `json:"fields,omitempty" yaml:"fields"`
	Uniq      string            `json:"uniq,omitempty" yaml:"uniq"`
	Relations map[string]*Query `json:"relations,omitempty" yaml:"relations"`
	With      map[string]*Query `json:"with,omitempty" yaml:"with"`
	Join      []Join            `json:"join,omitempty" yaml:"join"`
	Filter    []*Predicate      `json:"filter,omitempty" yaml:"filter"`
	Orders    []*Order          `json:"orders,omitempty" yaml:"orders"`
	Page      int64             `json:"page,omitempty" yaml:"page"`
	PerPage   int64             `json:"per_page,omitempty" yaml:"per_page"`
}

func NewQuery() *Query {
	return &Query{
		Relations: map[string]*Query{},
		With:      map[string]*Query{},
	}
}

func (q *Query) LimitOffset(bounded bool) (limit int, offset int) {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PerPage == 0 && bounded {
		q.PerPage = DefaultPerPage
	}
	if q.PerPage > MaxPerPage && bounded {
		q.PerPage = MaxPerPage
	}

	limit = int(q.PerPage)
	offset = int(q.PerPage * (q.Page - 1))
	return
}

func (q *Query) String() string {
	data, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(data)
}
