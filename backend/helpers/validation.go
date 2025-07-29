package helpers

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(data interface{}) (map[string]string, error) {
	err := validate.Struct(data)
	if err == nil {
		return nil, nil
	}

	fieldErrors := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		t := reflect.TypeOf(data)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		for _, fe := range ve {
			fieldName := fe.Field()
			if f, found := t.FieldByName(fieldName); found {
				jsonTag := strings.Split(f.Tag.Get("json"), ",")[0]
				if jsonTag != "" && jsonTag != "-" {
					fieldName = jsonTag
				}
			}
			fieldErrors[fieldName] = validationMessage(fe)
		}
	}

	log.Printf("Validation failed: %v", err)
	return fieldErrors, err
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Please enter a valid email address"
	case "min":
		return "Must be at least " + fe.Param() + " characters long"
	case "max":
		return "Must be at most " + fe.Param() + " characters long"
	case "gte":
		return "Must be greater than or equal to " + fe.Param()
	case "lte":
		return "Must be less than or equal to " + fe.Param()
	default:
		return "Invalid value"
	}
}
