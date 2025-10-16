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

type ValidatorError struct {
	Field string
	Err   error
}

var (
	ErrNoStruct   = errors.New("no struct")
	ErrFormatRule = errors.New("format rule")

	ErrValidationStringLength   = errors.New("string length")
	ErrValidationRegExpNotMatch = errors.New("regexp not match")
	ErrValidationIntNotMin      = errors.New("less min")
	ErrValidationIntNotMax      = errors.New("more max")
	ErrValidationNotIncludes    = errors.New("val not includes in list")
)

const (
	TypeValidationNested = "nested"
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	vErrors := make([]string, 0, len(v))
	for _, vError := range v {
		vErrors = append(vErrors, fmt.Sprintf("%s '%s'\n", vError.Field, vError.Err))
	}
	return strings.Join(vErrors, "")
}

type ValidatorErrors []ValidatorError

func (v ValidatorErrors) Error() string {
	vErrors := make([]string, 0, len(v))
	for _, vError := range v {
		vErrors = append(vErrors, fmt.Sprintf("%s '%s'\n", vError.Field, vError.Err))
	}
	return strings.Join(vErrors, "")
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	validatorEr := CheckRules(rv)
	if validatorEr != nil {
		return validatorEr
	}
	rt := rv.Type()
	validateErrors := make(ValidationErrors, 0, rt.NumField())
	if rt.Kind() != reflect.Struct {
		return append(validateErrors, ValidationError{Field: rt.Name(), Err: ErrNoStruct})
	}
	err := validateStruct(rv)
	if err != nil {
		validateErrors = append(validateErrors, err...)
	}
	if len(validateErrors) > 0 {
		return validateErrors
	}
	return nil
}

func CheckRules(rv reflect.Value) error {
	rt := rv.Type()
	errors := make(ValidatorErrors, 0, rt.NumField())
	if rt.Kind() != reflect.Struct {
		return append(errors, ValidatorError{Field: rt.Name(), Err: ErrFormatRule})
	}
	for i := 0; i < rv.NumField(); i++ {
		if !rt.Field(i).IsExported() {
			continue
		}
		field := rt.Field(i)
		tag := field.Tag

		if tag.Get("validate") == "" {
			continue
		}
		tagRules := strings.Split(tag.Get("validate"), "|")
		for _, rule := range tagRules {
			arRule := strings.Split(rule, ":")
			if arRule[0] == TypeValidationNested {
				err := CheckRules(rv.Field(i))
				if err != nil {
					errors = append(errors, ValidatorError{Field: field.Name, Err: err})
				}
				continue
			}
			typeRule := arRule[0]
			if len(arRule) > 2 {
				errors = append(errors, ValidatorError{Field: field.Name, Err: ErrFormatRule})
				continue
			}
			valRule := ""
			if len(arRule) == 2 {
				valRule = arRule[1]
			}
			err := checkTypeRules(field, typeRule, valRule)
			if err != nil {
				errors = append(errors, ValidatorError{Field: field.Name, Err: ErrFormatRule})
			}
		}
	}
	if len(errors) > 0 {
		return errors
	}

	return nil
}

func checkTypeRules(field reflect.StructField, typeRule string, valRule string) error {
	switch typeRule {
	case "in":
		if len(valRule) == 0 {
			return errors.New("empty val")
		}
		list := strings.Split(valRule, ",")
		for _, val := range list {
			if len(val) == 0 {
				return errors.New("empty val")
			}
			if field.Type.Kind() == reflect.Int {
				_, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
			}
		}
	case "len", "min", "max":
		if len(valRule) == 0 {
			return errors.New("empty val")
		}
		_, err := strconv.Atoi(valRule)
		return err
	case "regexp":
		if len(valRule) == 0 {
			return errors.New("empty val")
		}
		_, err := regexp.Compile(valRule)
		return err
	default:
		return errors.New("undefined type rule")
	}
	return nil
}

func extractValidateRules(field reflect.StructField) (rules map[string][]string, err error) {
	tag := field.Tag
	if tag.Get("validate") != "" {
		tagRules := strings.Split(tag.Get("validate"), "|")
		rules = make(map[string][]string, len(tagRules))
		for _, rule := range tagRules {
			arRule := strings.Split(rule, ":")
			switch {
			case len(arRule) == 2:
				rules[arRule[0]] = strings.Split(arRule[1], ",")
			case len(arRule) == 1 && arRule[0] == TypeValidationNested:
				rules[arRule[0]] = []string{}
			default:
				err = ErrFormatRule
			}
		}
	}
	return
}

func validateInt(fieldValue reflect.Value, ruleType string, ruleVal []string) error {
	val := fieldValue.Int()
	switch ruleType {
	case "min":
		ruleValInt, err := strconv.Atoi(ruleVal[0])
		if err != nil {
			return ErrFormatRule
		}
		if val < int64(ruleValInt) {
			return ErrValidationIntNotMin
		}
	case "max":
		ruleValInt, err := strconv.Atoi(ruleVal[0])
		if err != nil {
			return ErrFormatRule
		}
		if val > int64(ruleValInt) {
			return ErrValidationIntNotMax
		}
	case "in":
		ruleValInt := make([]int64, 0, len(ruleVal))
		for _, ruleValitem := range ruleVal {
			valint, err := strconv.Atoi(ruleValitem)
			if err != nil {
				return ErrFormatRule
			}
			ruleValInt = append(ruleValInt, int64(valint))
		}
		if slices.Index(ruleValInt, val) == -1 {
			return ErrValidationNotIncludes
		}
	}
	return nil
}

func validateString(fieldValue reflect.Value, ruleType string, ruleVal []string) error {
	val := fieldValue.String()
	switch ruleType {
	case "len":
		valint, err := strconv.Atoi(ruleVal[0])
		if err != nil {
			return ErrFormatRule
		}
		if len(val) != valint {
			return ErrValidationStringLength
		}
	case "regexp":
		regex, err := regexp.Compile(ruleVal[0])
		if err != nil {
			return err
		}
		if !regex.MatchString(val) {
			return ErrValidationRegExpNotMatch
		}
	case "in":
		if slices.Index(ruleVal, val) == -1 {
			return ErrValidationNotIncludes
		}
	}
	return nil
}

func validateStruct(rv reflect.Value) ValidationErrors {
	rt := rv.Type()
	validateErrors := make(ValidationErrors, 0, rt.NumField())
	for i := 0; i < rv.NumField(); i++ {
		if !rt.Field(i).IsExported() {
			continue
		}
		rules, errRules := extractValidateRules(rt.Field(i))
		if errRules != nil {
			validateErrors = append(validateErrors, ValidationError{Field: rt.Field(i).Name, Err: errRules})
		}
		field := rv.Field(i)
		for ruleType, ruleVal := range rules {
			var errValidate error
			//nolint:exhaustive
			switch field.Kind() {
			case reflect.Struct:
				if ruleType == TypeValidationNested {
					err := validateStruct(field)
					if len(err) > 0 {
						validateErrors = append(validateErrors, ValidationError{Field: rt.Field(i).Name, Err: err})
					}
				}
			case reflect.Int, reflect.String, reflect.Slice:
				errValidate = validateItem(field, ruleType, ruleVal)
			}
			if errValidate != nil {
				validateErrors = append(validateErrors, ValidationError{Field: rt.Field(i).Name, Err: errValidate})
			}
		}
	}
	if len(validateErrors) > 0 {
		return validateErrors
	}
	return nil
}

func validateItem(field reflect.Value, ruleType string, ruleVal []string) error {
	//nolint:exhaustive
	switch field.Kind() {
	case reflect.Slice:
		for si := 0; si < field.Len(); si++ {
			err := validateItem(field.Index(si), ruleType, ruleVal)
			if err != nil {
				return err
			}
		}
	case reflect.Int:
		err := validateInt(field, ruleType, ruleVal)
		if err != nil {
			return err
		}
	case reflect.String:
		err := validateString(field, ruleType, ruleVal)
		if err != nil {
			return err
		}
	}
	return nil
}
