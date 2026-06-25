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

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder

	for i, vv := range v {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(vv.Error())
	}

	return sb.String()
}

type Rule struct {
	Name string
	Raw  string
}

func parseRules(validateTag string) ([]Rule, error) {
	splitsTag := strings.Split(validateTag, "|")
	rules := make([]Rule, 0, len(splitsTag))

	for _, part := range splitsTag {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		nameRaw := strings.SplitN(part, ":", 2)
		if len(nameRaw) != 2 {
			return nil, fmt.Errorf("invalid validate tag %q", part)
		}
		rules = append(rules, Rule{
			Name: nameRaw[0],
			Raw:  nameRaw[1],
		})
	}

	return rules, nil
}

func Validate(v interface{}) error {
	reflectVal := reflect.ValueOf(v)
	if reflectVal.Kind() == reflect.Ptr {
		if reflectVal.IsNil() {
			return errors.New("nil pointer")
		}
		reflectVal = reflectVal.Elem()
	}

	if reflectVal.Kind() != reflect.Struct {
		return errors.New("v must be a struct")
	}

	var result ValidationErrors
	for i := 0; i < reflectVal.NumField(); i++ {
		field := reflectVal.Type().Field(i)
		fieldVal := reflectVal.Field(i)
		validationErrors, err := validateField(field, fieldVal)
		if err != nil {
			return err
		}
		result = append(result, validationErrors...)
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

func validateField(field reflect.StructField, fieldValue reflect.Value) (ValidationErrors, error) {
	if !field.IsExported() {
		return nil, nil
	}

	validateTag, exists := field.Tag.Lookup("validate")
	if !exists {
		return nil, nil
	}

	rules, err := parseRules(validateTag)
	if err != nil {
		return nil, err
	}

	//nolint:exhaustive
	switch field.Type.Kind() {
	case reflect.Int:
		return validateSingleIntField(fieldValue.Int(), field.Name, rules)
	case reflect.String:
		return validateSingleStringField(fieldValue.String(), field.Name, rules)
	case reflect.Slice:
		return validateSliceField(field.Name, fieldValue, rules)
	default:
		return nil, fmt.Errorf("unsupported type %s", field.Type)
	}
}

func validateSliceField(fieldName string, sliceVal reflect.Value, rules []Rule) (ValidationErrors, error) {
	var errs ValidationErrors

	for i := 0; i < sliceVal.Len(); i++ {
		var validationFieldErrors ValidationErrors
		var err error

		elem := sliceVal.Index(i)
		elemKind := sliceVal.Type().Elem().Kind()
		elemField := fmt.Sprintf("%s[%d]", fieldName, i)

		//nolint:exhaustive
		switch elemKind {
		case reflect.String:
			validationFieldErrors, err = validateSingleStringField(elem.String(), elemField, rules)
			if err != nil {
				return nil, err
			}
		case reflect.Int:
			validationFieldErrors, err = validateSingleIntField(elem.Int(), elemField, rules)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported slice element type %s", sliceVal.Type().Elem())
		}

		if validationFieldErrors != nil {
			errs = append(errs, validationFieldErrors...)
		}
	}

	return errs, nil
}

func validateSingleStringField(value, fieldName string, rules []Rule) (ValidationErrors, error) {
	if len(rules) == 0 {
		return nil, nil
	}

	var errs ValidationErrors
	for _, rule := range rules {
		switch rule.Name {
		case "len":
			lenValue, err := strconv.Atoi(rule.Raw)
			if err != nil {
				return nil, err
			}
			if len(value) != lenValue {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("expected length %v, got %v", lenValue, len(value)),
				})
			}
		case "regexp":
			re, err := regexp.Compile(rule.Raw)
			if err != nil {
				return nil, err
			}
			if !re.MatchString(value) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("regexp does not match %s", rule.Raw),
				})
			}
		case "in":
			var find bool
			for _, part := range strings.Split(rule.Raw, ",") {
				if strings.TrimSpace(part) == value {
					find = true
				}
			}

			if !find {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("%s is not in %s", rule.Name, rule.Raw),
				})
			}
		default:
			return nil, fmt.Errorf("unsupported rule %q", rule.Name)
		}
	}

	return errs, nil
}

func validateSingleIntField(value int64, fieldName string, rules []Rule) (ValidationErrors, error) {
	if len(rules) == 0 {
		return nil, nil
	}

	var errs ValidationErrors
	for _, rule := range rules {
		switch rule.Name {
		case "min":
			minValue, err := strconv.Atoi(rule.Raw)
			if err != nil {
				return nil, err
			}
			if value < int64(minValue) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value must be greater than %d", minValue),
				})
			}
		case "max":
			maxValue, err := strconv.Atoi(rule.Raw)
			if err != nil {
				return nil, err
			}
			if value > int64(maxValue) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value must be less than %d", maxValue),
				})
			}
		case "in":
			contains, err := containsValue(rule.Raw, value)
			if err != nil {
				return nil, err
			}

			if !contains {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value must contain %s", rule.Raw),
				})
			}
		default:
			return nil, fmt.Errorf("unsupported rule %q", rule.Name)
		}
	}

	return errs, nil
}

func containsValue(s string, v int64) (bool, error) {
	for _, part := range strings.Split(s, ",") {
		intValue, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return false, err
		}

		if int64(intValue) == v {
			return true, nil
		}
	}

	return false, nil
}
