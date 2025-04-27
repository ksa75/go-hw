package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
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

		// if validateTag == "nested" {
		// 	if field.Kind() == reflect.Struct {
		// 		if err := Validate(field.Interface()); err != nil {
		// 			if verrs, ok := err.(ValidationErrors); ok {
		// 				for _, verr := range verrs {
		// 					validationErrors = append(validationErrors, ValidationError{
		// 						Field: fieldType.Name + "." + verr.Field,
		// 						Err:   verr.Err,
		// 					})
		// 				}
		// 			} else {
		// 				return err
		// 			}
		// 		}
		// 	}
		// 	continue
		// }

		rules := strings.Split(validateTag, "|")

		switch field.Kind() {
		case reflect.String:
			for _, rule := range rules {
				if err := validateString(field.String(), rule); err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   err,
					})
				}
			}

		case reflect.Int:
			for _, rule := range rules {
				if err := validateInt(int(field.Int()), rule); err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   err,
					})
				}
			}

		case reflect.Slice:
			switch field.Type().Elem().Kind() {
			case reflect.String:
				for idx := 0; idx < field.Len(); idx++ {
					for _, rule := range rules {
						if err := validateString(field.Index(idx).String(), rule); err != nil {
							validationErrors = append(validationErrors, ValidationError{
								Field: fmt.Sprintf("%s[%d]", fieldType.Name, idx),
								Err:   err,
							})
						}
					}
				}
			case reflect.Int:
				for idx := 0; idx < field.Len(); idx++ {
					for _, rule := range rules {
						if err := validateInt(int(field.Index(idx).Int()), rule); err != nil {
							validationErrors = append(validationErrors, ValidationError{
								Field: fmt.Sprintf("%s[%d]", fieldType.Name, idx),
								Err:   err,
							})
						}
					}
				}
			}
		}
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
		found := false
		for _, v := range values {
			if value == v {
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
