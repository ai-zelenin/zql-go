package zql

type MutatorPredicate interface {
	Mutate(p *Predicate) error
}

type MutatePredicateFunc func(p *Predicate) error

func (v MutatePredicateFunc) Mutate(p *Predicate) error {
	return v(p)
}
