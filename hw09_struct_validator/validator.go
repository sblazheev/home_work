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

var (
	ErrNoStruct = errors.New("no struct")
	ErrNoParts  = errors.New("no name and value in tag")

	ErrValidationRule              = errors.New("format rules")
	ErrValidationStringLength      = errors.New("string length")
	ErrValidationRegExpNotMatch    = errors.New("regexp not match")
	ErrValidationNotIncludesString = errors.New("string not includes in list")
	ErrValidationIntNotMin         = errors.New("less min")
	ErrValidationIntNotMax         = errors.New("more max")
	ErrValidationNotIncludes       = errors.New("val not includes in list")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errors := make([]string, 0, len(v))
	for _, vError := range v {
		errors = append(errors, fmt.Sprintf("%s '%s'\n", vError.Field, vError.Err))
	}
	return strings.Join(errors, "")
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
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
			case len(arRule) == 1 && arRule[0] == "nested":
				rules[arRule[0]] = []string{}
			default:
				err = ErrValidationRule
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
			return ErrValidationRule
		}
		if val < int64(ruleValInt) {
			return ErrValidationIntNotMin
		}
	case "max":
		ruleValInt, err := strconv.Atoi(ruleVal[0])
		if err != nil {
			return ErrValidationRule
		}
		if val > int64(ruleValInt) {
			return ErrValidationIntNotMax
		}
	case "in":
		ruleValInt := make([]int64, 0, len(ruleVal))
		for _, ruleValitem := range ruleVal {
			valint, err := strconv.Atoi(ruleValitem)
			if err != nil {
				return ErrValidationRule
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
			return ErrValidationRule
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
				if ruleType == "nested" {
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
