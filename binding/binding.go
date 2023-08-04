package binding

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type FieldErrors = map[string][]string

func schemaErrorToText(err error) string {
	switch err := err.(type) {
	case schema.ConversionError:
		switch err.Type.Name() {
		case "int":
			return "Not a valid integer"
		}
	case schema.EmptyFieldError:
		return "Required field"
	}
	return "Invalid"
}

func schemaMultiErrorToFieldErrors(err schema.MultiError) FieldErrors {
	errorInfo := make(map[string][]string)
	for field, fieldError := range err {
		if errorInfo[field] == nil {
			errorInfo[field] = make([]string, 0, 1)
		}
		errorInfo[field] = append(errorInfo[field], schemaErrorToText(fieldError))
	}
	return errorInfo
}

func validationErrorToText(err validator.FieldError) string {
	switch err.Tag() {
	case "min":
		return fmt.Sprintf("must be >= %s", err.Param())
	case "max":
		return fmt.Sprintf("must be <= %s", err.Param())
	case "url":
		return "Invalid url"
	}
	return "Invalid"
}

func validationErrorsToFieldErrors(err validator.ValidationErrors) FieldErrors {
	errorInfo := make(map[string][]string)
	for _, fieldError := range err {
		field := fieldError.Field()
		if errorInfo[field] == nil {
			errorInfo[field] = make([]string, 0, 1)
		}
		errorInfo[field] = append(errorInfo[field], validationErrorToText(fieldError))
	}
	return errorInfo
}

func createValidate() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("schema"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return validate
}

var decoder = schema.NewDecoder()
var validate = createValidate()

var genericError = map[string][]string{
	"error": {"An error has occurred"},
}

func ShouldBind(dst any, src map[string][]string) FieldErrors {
	fieldErrors := make(map[string][]string)
	if err := decoder.Decode(dst, src); err != nil {
		err, ok := err.(schema.MultiError)
		if err != nil && !ok {
			return genericError
		}
		fieldErrors = schemaMultiErrorToFieldErrors(err)
	}
	if err := validate.Struct(dst); err != nil {
		err, ok := err.(validator.ValidationErrors)
		if err != nil && !ok {
			log.Print(err)
			return genericError
		}
		errs := validationErrorsToFieldErrors(err)
		// Only merge validation errors into the result if no decoding errors occurred
		for field, errorTexts := range errs {
			if fieldErrors[field] == nil {
				fieldErrors[field] = errorTexts
			}
		}
	}
	if len(fieldErrors) == 0 {
		return nil
	}
	return fieldErrors
}
