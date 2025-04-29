package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, ve := range v {
		sb.WriteString(fmt.Sprintf("Field '%s': %s\n", ve.Field, ve.Err.Error()))
	}
	return sb.String()
}

func Validate(v any) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return errors.New("only structs are supported")
	}

	var validationErrors ValidationErrors

	typ := val.Type()
	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)
		validateTag := fieldType.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, "|")
		kind := field.Kind()

		if kind == reflect.String {
			for _, rule := range rules {
				if err := validateString(field.String(), rule); err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   err,
					})
				}
			}
			continue
		}

		if kind == reflect.Int {
			for _, rule := range rules {
				if err := validateInt(int(field.Int()), rule); err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   err,
					})
				}
			}
			continue
		}

		if kind == reflect.Slice {
			validateSliceField(field, fieldType, rules, &validationErrors)
			continue
		}
	}
	if isRegexError(validationErrors) {
		return errors.New("invalid regexp")
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateString(value string, rule string) error {
	switch {
	case strings.HasPrefix(rule, "len:"):
		expectedLen, _ := strconv.Atoi(strings.TrimPrefix(rule, "len:"))
		if len(value) != expectedLen {
			return fmt.Errorf("length must be %d", expectedLen)
		}

	case strings.HasPrefix(rule, "regexp:"):
		pattern := strings.TrimPrefix(rule, "regexp:")
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regexp: %s", pattern)
		}
		if !re.MatchString(value) {
			return fmt.Errorf("must match regexp %s", pattern)
		}

	case strings.HasPrefix(rule, "in:"):
		values := strings.Split(strings.TrimPrefix(rule, "in:"), ",")
		found := slices.Contains(values, value)
		if !found {
			return fmt.Errorf("must be one of [%s]", strings.Join(values, ", "))
		}
	}
	return nil
}

func validateInt(value int, rule string) error {
	switch {
	case strings.HasPrefix(rule, "min:"):
		minValue, _ := strconv.Atoi(strings.TrimPrefix(rule, "min:"))
		if value < minValue {
			return fmt.Errorf("must be >= %d", minValue)
		}

	case strings.HasPrefix(rule, "max:"):
		maxValue, _ := strconv.Atoi(strings.TrimPrefix(rule, "max:"))
		if value > maxValue {
			return fmt.Errorf("must be <= %d", maxValue)
		}

	case strings.HasPrefix(rule, "in:"):
		raw := strings.TrimPrefix(rule, "in:")
		values := strings.Split(raw, ",")
		found := false
		for _, v := range values {
			intV, _ := strconv.Atoi(v)
			if value == intV {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("must be one of [%s]", strings.Join(values, ", "))
		}
	}
	return nil
}

func validateSliceField(fld reflect.Value, fldType reflect.StructField, rules []string, valErrors *ValidationErrors) {
	elemKind := fld.Type().Elem().Kind()

	if elemKind == reflect.String {
		for idx := range fld.Len() {
			for _, rule := range rules {
				if err := validateString(fld.Index(idx).String(), rule); err != nil {
					*valErrors = append(*valErrors, ValidationError{
						Field: fmt.Sprintf("%s[%d]", fldType.Name, idx),
						Err:   err,
					})
				}
			}
		}
		return
	}

	if elemKind == reflect.Int {
		for idx := range fld.Len() {
			for _, rule := range rules {
				if err := validateInt(int(fld.Index(idx).Int()), rule); err != nil {
					*valErrors = append(*valErrors, ValidationError{
						Field: fmt.Sprintf("%s[%d]", fldType.Name, idx),
						Err:   err,
					})
				}
			}
		}
		return
	}
}

func isRegexError(actualErr error) bool {
	return strings.Contains(actualErr.Error(), "invalid regexp")
}
