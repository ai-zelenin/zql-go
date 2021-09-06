package dao

type SyntaxV1 struct {
	qb *QueryBuilder
}

func NewSyntaxV1() *SyntaxV1 {
	return &SyntaxV1{qb: NewQueryBuilder()}
}

func (s *SyntaxV1) JoinAnd(p1, p2 *Predicate) *Predicate {
	if p1.Op == AND {
		predicates := p1.Value.([]*Predicate)
		predicates = append(predicates, p2)
		p1.Value = predicates
		return p1
	}
	return And(p1, p2)
}

func (s *SyntaxV1) JoinOr(p1, p2 *Predicate) *Predicate {
	if p1.Op == OR {
		predicates := p1.Value.([]*Predicate)
		predicates = append(predicates, p2)
		p1.Value = predicates
		return p1
	}
	return Or(p1, p2)
}

func (s *SyntaxV1) Env() map[string]interface{} {
	return map[string]interface{}{
		"qb":              s.qb,
		"asc":             Asc,
		"desc":            Desc,
		"AND":             And,
		"OR":              Or,
		"EqOperator":      NewOperator("==", Eq),
		"NeqOperator":     NewOperator("!=", Neq),
		"GtOperator":      NewOperator(">", Gt),
		"GteOperator":     NewOperator(">=", Gte),
		"LtOperator":      NewOperator("<", Lt),
		"LteOperator":     NewOperator("<=", Lte),
		"In":              NewOperator("in", In),
		"JoinAndOperator": NewOperator("&&", s.JoinAnd),
		"JoinOrOperator":  NewOperator("||", s.JoinOr),
	}
}
