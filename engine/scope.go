package engine

type Scope interface {
	Get(k string) (v interface{}, ok bool)
	Set(k string, v interface{})
}

type MapScope struct {
	scope map[string]any
}

func NewMapScope(m map[string]any) *MapScope {
	if m == nil {
		m = make(map[string]any)
	}
	return &MapScope{
		scope: m,
	}
}

func (m *MapScope) Get(k string) (v interface{}, ok bool) {
	v, ok = m.scope[k]
	return v, ok
}

func (m *MapScope) Set(k string, v interface{}) {
	m.scope[k] = v
}
