package zql

type ValidatorConfig struct {
	CheckAcceptableFields bool
	AcceptableFields      map[string]bool

	CheckFieldValueTypeMap bool
	FieldValueTypeMap      map[string]string

	CheckOpValueTypeMap bool
	OpValueTypeMap      map[string]string
}
