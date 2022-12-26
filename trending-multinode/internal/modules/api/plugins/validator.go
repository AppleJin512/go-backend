package plugins

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"reflect"
)

type Validator struct {
	goPlaygroundValidate *validator.Validate
}

func NewValidator() *Validator {
	goPlaygroundValidate := validator.New()
	goPlaygroundValidate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	customValidators := map[string]validator.FuncCtx{}
	for tag, fun := range customValidators {
		if err := goPlaygroundValidate.RegisterValidationCtx(tag, fun); err != nil {
			logrus.Errorf("register validation ctx error: %v", err)
		}
	}

	return &Validator{
		goPlaygroundValidate: goPlaygroundValidate,
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.goPlaygroundValidate.Struct(i)
}
