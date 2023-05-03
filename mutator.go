package zql

type Mutator interface {
	Mutate(q *Query) error
}

type ExtendableMutator struct {
	Mutators          []Mutator
	PredicateMutators []MutatorPredicate
}

func NewExtendableMutator() *ExtendableMutator {
	return &ExtendableMutator{
		Mutators:          make([]Mutator, 0),
		PredicateMutators: make([]MutatorPredicate, 0),
	}
}

func (e *ExtendableMutator) AddMutator(v Mutator) {
	e.Mutators = append(e.Mutators, v)
}

func (e *ExtendableMutator) AddPredicateMutator(v MutatorPredicate) {
	e.PredicateMutators = append(e.PredicateMutators, v)
}

func (e *ExtendableMutator) Mutate(q *Query) error {
	err := e.mutateFilter(q)
	if err != nil {
		return err
	}
	for _, mutator := range e.Mutators {
		err = mutator.Mutate(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ExtendableMutator) mutateFilter(q *Query) error {
	for _, predicate := range q.Filter {
		err := predicate.Walk(func(parent Node, current Node, lvl int) error {
			currentPr := current.(*Predicate)
			err := e.mutatePredicate(currentPr)
			if err != nil {
				return err
			}
			return nil
		}, nil, 0)
		if err != nil {
			return err
		}
	}
	for _, query := range q.Relations {
		err := e.mutateFilter(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ExtendableMutator) mutatePredicate(p *Predicate) error {
	for _, mutator := range e.PredicateMutators {
		if mutator != nil {
			err := mutator.Mutate(p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
