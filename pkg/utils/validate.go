package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Validate(value any) []string {
	errors := make([]string, 0)
	v := reflect.ValueOf(value)

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		if field.Kind() == reflect.Struct {
			errors = append(errors, Validate(field.Interface())...)
			continue
		}

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		errors = append(errors, validateField(field, fieldType.Name, tag)...)
	}

	return errors
}

func validateField(val reflect.Value, fieldName string, tag string) []string {
	errors := make([]string, 0)
	validations := strings.Split(tag, ",")

	for _, validation := range validations {
		if strings.HasPrefix(validation, "min:") {
			if errorString := validateMin(validation, val, fieldName); errorString != "" {
				errors = append(errors, errorString)
			}
		} else if strings.HasPrefix(validation, "max:") {
			if errorString := validateMax(validation, val, fieldName); errorString != "" {
				errors = append(errors, errorString)
			}
		} else if validation == "required" {
			if errorString := validateRequired(val, fieldName); errorString != "" {
				errors = append(errors, errorString)
			}
		} else if validation == "password" {
			if errorStrings := validatePassword(val, fieldName); len(errorStrings) > 0 {
				errors = append(errors, errorStrings...)
			}
		} else if validation == "email" {
			if errorString := validateEmail(val, fieldName); errorString != "" {
				errors = append(errors, errorString)
			}
		}
	}

	return errors
}

func validateMin(validate string, val reflect.Value, fieldName string) string {
	min, _ := strconv.Atoi(strings.TrimPrefix(validate, "min:"))
	if len(val.String()) < min {
		return fmt.Sprintf("Field '%s' must be at least %d characters long", fieldName, min)
	}
	return ""
}

func validateMax(validate string, val reflect.Value, field string) string {
	max, _ := strconv.Atoi(strings.TrimPrefix(validate, "max:"))
	if len(val.String()) > max {
		return fmt.Sprintf("Field '%s' must be at most %d characters long", field, max)
	}
	return ""
}

func validateRequired(val reflect.Value, fieldName string) string {
	if val.String() == "" {
		return fmt.Sprintf("Field '%s' is required", fieldName)
	}
	return ""
}

func validateEmail(val reflect.Value, fieldName string) string {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if ok := regexp.MustCompile(regex).MatchString(val.String()); !ok {
		return fmt.Sprintf("Field '%s' is not a valid email", fieldName)
	}
	return ""
}

func validatePassword(val reflect.Value, fieldName string) []string {
	password := val.String()
	errors := make([]string, 0)

	number := false
	specialCharacter := false
	upperCase := false
	lowerCase := false

	for _, c := range []byte(password) {
		if IsNumber(c) {
			number = true
		} else if IsSpecialCharacter(c) {
			specialCharacter = true
		} else if c >= 'a' && c <= 'z' {
			lowerCase = true
		} else if c >= 'A' && c <= 'Z' {
			upperCase = true
		}
	}

	if !number {
		errors = append(errors, fmt.Sprintf("Field '%s' requires at least one number", fieldName))
	}
	if !specialCharacter {
		errors = append(errors, fmt.Sprintf("Field '%s' requires at least one special character", fieldName))
	}
	if !upperCase {
		errors = append(errors, fmt.Sprintf("Field '%s' requires at least one upper case letter", fieldName))
	}
	if !lowerCase {
		errors = append(errors, fmt.Sprintf("Field '%s' requires at least one lower case letter", fieldName))
	}
	// if len(password) < 8 || len(password) > 32 {
	// 	errors = append(errors, "Password length is not between 8 and 32 characters long")
	// }

	return errors
}

func IsNumber(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}

	return false
}

func IsSpecialCharacter(c byte) bool {
	if (c >= '!' && c <= '/') ||
		(c >= ':' && c <= '@') ||
		(c >= '[' && c <= '^') ||
		(c >= '{' && c <= '~') {
		return true
	}

	return false
}
