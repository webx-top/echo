package echo

// Validator is the interface that wraps the Validate method.
type Validator interface {
	Validate(i ...interface{}) error
}

var DefaultNopValidate Validator = &NopValidation{}

type NopValidation struct {
}

func (v *NopValidation) Validate(_ ...interface{}) error {
	return nil
}
