package compare

type UntypedComparator struct {
	*TypedComparator
}

func NewUntypedComparator() *UntypedComparator {
	return &UntypedComparator{
		TypedComparator: NewTypedComparator(),
	}
}

func (u *UntypedComparator) Eq(v1, v2 any) (result bool) {
	result, err := u.TypedComparator.Eq(v1, v2)
	if err != nil {
		return false
	}
	return
}

func (u *UntypedComparator) Neq(v1, v2 any) (result bool) {
	result, err := u.TypedComparator.Neq(v1, v2)
	if err != nil {
		return false
	}
	return
}

func (u *UntypedComparator) Gt(v1, v2 any) (result bool) {
	result, err := u.TypedComparator.Gt(v1, v2)
	if err != nil {
		return false
	}
	return
}

func (u *UntypedComparator) Gte(v1, v2 any) (result bool) {
	result, err := u.TypedComparator.Gte(v1, v2)
	if err != nil {
		return false
	}
	return
}

func (u *UntypedComparator) Lt(v1, v2 any) (result bool) {
	result, err := u.TypedComparator.Lt(v1, v2)
	if err != nil {
		return false
	}
	return
}

func (u *UntypedComparator) Lte(v1, v2 any) (result bool) {
	result, err := u.TypedComparator.Lte(v1, v2)
	if err != nil {
		return false
	}
	return
}
