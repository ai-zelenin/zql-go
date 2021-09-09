package zql

const DefaultMaxPredicateNumber = 20

type ValidatorConfig struct {
	Validators          []Validator
	PredicateValidators []ValidatorPredicate
	MaxPredicateNumber  int
}

func NewValidatorConfigForModel(model interface{}, tagName string) *ValidatorConfig {
	fieldsValidator := NewValidatorPredicateFields()
	valueValidator := NewValidatorPredicateValues()
	opsValidator := NewValidatorPredicateOps()

	fields := fieldsFromModel(model, tagName)
	for fieldName, valueType := range fields {
		fieldsValidator.AddField(fieldName)
		valueValidator.AddFieldValuePair(fieldName, valueType)
	}

	return &ValidatorConfig{
		PredicateValidators: []ValidatorPredicate{
			fieldsValidator,
			valueValidator,
			opsValidator,
		},
		MaxPredicateNumber: DefaultMaxPredicateNumber,
	}
}
