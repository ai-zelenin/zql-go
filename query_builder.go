package zql

type QueryBuilder struct {
	query  *Query
	relMap map[string]*QueryBuilder
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		query:  NewQuery(),
		relMap: map[string]*QueryBuilder{},
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

func (qb *QueryBuilder) Relation(rel string) *QueryBuilder {
	newQb := NewQueryBuilder()
	qb.relMap[rel] = newQb
	return newQb
}

func (qb *QueryBuilder) Order(order ...*Order) *QueryBuilder {
	qb.query.Orders = append(qb.query.Orders, order...)
	return qb
}

func (qb *QueryBuilder) Uniq(field string) {
	qb.query.Uniq = field
}

func (qb *QueryBuilder) Build() *Query {
	q := qb.query
	for rel, builder := range qb.relMap {
		q.Relations[rel] = builder.Build()
	}
	return q
}
