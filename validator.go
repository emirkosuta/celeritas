package celeritas

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Validation struct {
	Errors map[string]string
	V      *validator.Validate
}

func (c *Celeritas) Validator() *Validation {
	return &Validation{
		Errors: make(map[string]string),
		V:      validate,
	}
}

func (v *Validation) ValidateStruct(s interface{}) *Validation {
	err := v.V.Struct(s)
	if err != nil {
		// If the type of the error is validator.ValidationErrors, then we can
		// range over the individual field errors.
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range errs {
				v.AddError(err.Field(), err.Error())
			}
		}
	}
	return v
}

func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validation) AddError(key string, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}
