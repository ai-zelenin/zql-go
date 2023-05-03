package engine

type ExtendableEngine struct {
	Functions map[string]ExecuteFunc
	cfg       *Config
}

func NewExtendableEngine(cfg *Config) *ExtendableEngine {
	if cfg == nil {
		cfg = DefaultConfig
	}
	if cfg.Comparator == nil {
		cfg.Comparator = DefaultConfig.Comparator
	}
	if cfg.CompileFunc == nil {
		cfg.CompileFunc = DefaultConfig.CompileFunc
	}
	c := cfg.Comparator
	funcs := map[string]ExecuteFunc{
		"and": calcAND,
		"or":  calcOR,
		"eq":  newResolveAndCompareFunc(c.Eq),
		"neq": newResolveAndCompareFunc(c.Neq),
		"gt":  newResolveAndCompareFunc(c.Gt),
		"gte": newResolveAndCompareFunc(c.Gte),
		"lt":  newResolveAndCompareFunc(c.Lt),
		"lte": newResolveAndCompareFunc(c.Lte),
	}
	for k, executeFunc := range cfg.PredicateFunctions {
		funcs[k] = executeFunc
	}
	return &ExtendableEngine{
		cfg:       cfg,
		Functions: funcs,
	}
}

func (e *ExtendableEngine) Compile(predicate Predicate) (ep *ExecutablePredicate, err error) {
	return e.cfg.CompileFunc(e.Functions, predicate, e.cfg.PreHooks, e.cfg.PostHooks)
}

func calcAND(e *ExecutablePredicate, s Scope) (result bool, err error) {
	predicateList, ok := e.Value.([]*ExecutablePredicate)
	if !ok {
		return false, ErrBadValueType
	}
	for _, predicate := range predicateList {
		result, err = predicate.Exec(s)
		if err != nil {
			return false, err
		}
		if result == false {
			return result, nil
		}
	}
	return true, nil
}

func calcOR(e *ExecutablePredicate, s Scope) (result bool, err error) {
	predicateList, ok := e.Value.([]*ExecutablePredicate)
	if !ok {
		return false, ErrBadValueType
	}
	for _, predicate := range predicateList {
		result, err = predicate.Exec(s)
		if err != nil {
			return false, err
		}
		if result == true {
			return result, nil
		}
	}
	return false, nil
}

func newResolveAndCompareFunc(cmpFunc func(v1, v2 any) (result bool, err error)) ExecuteFunc {
	return func(e *ExecutablePredicate, scope Scope) (result bool, err error) {
		v1, ok := scope.Get(e.Field)
		if !ok {
			return false, ErrVariableNotFound
		}
		return cmpFunc(v1, e.Value)
	}
}
