package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")

		if name == "-" {
			return ""
		}

		return name
	})

	return validate
}
