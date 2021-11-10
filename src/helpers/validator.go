package helpers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateInput(dataSet interface{}) (bool, map[string]string) {
	var validate *validator.Validate
	validate = validator.New()

	err := validate.Struct(dataSet)

	if err != nil {
		// validate syntax error
		if err, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		// validation errors occurred
		errors := make(map[string]string)

		reflected := reflect.ValueOf(dataSet)

		for _, err := range err.(validator.ValidationErrors) {
			// find field by by name and get json value
			field, _ := reflected.Type().FieldByName(err.StructField())

			var name string

			if name = field.Tag.Get("json"); name == "" {
				name = strings.ToLower(err.StructField())
			}

			fmt.Println(err.Tag())

			switch err.Tag() {
			case "required":
				errors[name] = "The " + name + " field is required"

			case "email":
				errors[name] = "The " + name + " field must be a valid email"

			// case "gte":
			// 	errors[name] = name + " must be greater than 7 charactors"

			default:
				errors[name] = "The " + name + " field is invalid"
			}

		}
		return false, errors
	}

	return true, nil
}
