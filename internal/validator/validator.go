// Package validator provides utilities and helpers for validating form data and inputs.
package validator

import (
	"strings"
	"unicode/utf8"
)

// Validator contains maps and slices of validation errors.
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid returns true if the Validator contains no errors.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError adds an error message for a specific field, if one does not already exist.
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// AddNonFieldError adds an error message not tied to any specific field.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField adds an error message if the validation check is not ok.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank returns true if a value is not empty and not only whitespace.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxRunes returns true if a value contains no more than n characters (runes).
func MaxRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt returns true if a value is present in a permitted list of integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for _, v := range permittedValues {
		if value == v {
			return true
		}
	}
	return false
}
