package dao

type QueryBuilder struct {
	query *Query
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		query: NewQuery(),
	}
}

func (qb *QueryBuilder) Filter(p ...*Predicate) *QueryBuilder {
	qb.query.Filter = append(qb.query.Filter, p...)
	return qb
}

func (qb *QueryBuilder) Page(page int64, perPage ...int64) *QueryBuilder {
	qb.query.Page = page
	if len(perPage) > 0 {
		qb.query.PerPage = perPage[0]
	}
	return qb
}

func (qb *QueryBuilder) Relation(rels ...string) *QueryBuilder {
	qb.query.Relations = append(qb.query.Relations, rels...)
	return qb
}

func (qb *QueryBuilder) Order(order ...*Order) *QueryBuilder {
	qb.query.Orders = append(qb.query.Orders, order...)
	return qb
}

func (qb *QueryBuilder) Distinct(field string) {
	qb.query.Distinct = field
}

func (qb *QueryBuilder) Build() *Query {
	return qb.query
}
