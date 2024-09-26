package exp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

const (
	invalidBodyText  = "invalid body"
	badRequestStatus = http.StatusBadRequest

	unprocessableEntityStatus = http.StatusUnprocessableEntity
)

func IsValidationExp(err error) *ValidationExp {
	if err == nil {
		return nil
	}

	var validationExp *ValidationExp
	if errors.As(err, &validationExp) {
		return validationExp
	}
	return nil
}

func isValidationErrors(err error) validator.ValidationErrors {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return validationErrors
	}
	return nil
}

func isUnmarshalTypeError(err error) *json.UnmarshalTypeError {
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return unmarshalTypeError
	}
	return nil
}

func isSyntaxError(err error) *json.SyntaxError {
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return syntaxError
	}
	return nil
}

func isInvalidBody(err error) bool {
	if err == nil {
		return false
	} else if err.Error() == "EOF" {
		return true
	}
	return isUnmarshalTypeError(err) != nil || isSyntaxError(err) != nil
}

type ValidationExp struct {
	Status int      `json:"status"`
	Errors []string `json:"errors"`
}

func (e *ValidationExp) Error() string {
	return fmt.Sprintf("Status: %d, Errors: %s", e.Status, e.Errors)
}

func NewValidationExp(e error) *ValidationExp {
	if e == nil {
		return nil
	}

	validationErrors := isValidationErrors(e)
	if validationErrors == nil {
		if ok := isInvalidBody(e); ok {
			return &ValidationExp{
				Status: badRequestStatus,
				Errors: []string{invalidBodyText},
			}
		}
		return nil
	}
	return &ValidationExp{
		Status: unprocessableEntityStatus,
		Errors: parseValidationErrors(validationErrors),
	}
}

func parseValidationErrors(validationErrors validator.ValidationErrors) []string {
	var errs []string
	for _, err := range validationErrors {
		fieldName := strings.ToLower(err.Field()) // example: "username"
		actualTag := err.ActualTag()              // example: "min"
		param := err.Param()                      // example: "6"

		// Generating an understandable error message
		message := generateErrorMessage(fieldName, actualTag, param)
		errs = append(errs, message)
	}
	return errs
}

func generateErrorMessage(fieldName, actualTag, param string) string {
	switch actualTag {
	case "required":
		return fmt.Sprintf("Field '%s' is required.", fieldName)
	case "email":
		return fmt.Sprintf("Field '%s' is not a valid email address.", fieldName)
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s characters long.", fieldName, param)
	case "max":
		return fmt.Sprintf("Field '%s' must be no more than %s characters long.", fieldName, param)
	case "ascii":
		return fmt.Sprintf("Field '%s' must contain only ASCII characters.", fieldName)
	case "bool":
		return fmt.Sprintf("Field '%s' must be a boolean.", fieldName)
	default:
		return fmt.Sprintf("Field '%s' failed validation due to: %s", fieldName, actualTag)
	}
}
