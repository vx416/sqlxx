package parsers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vicxu416/sqlxx/errors"
)

func ParseConditions(source interface{}, cond map[string]interface{}, allowEmpty bool) error {
	val := reflect.ValueOf(source)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			fieldName := val.Type().Field(i).Tag.Get(Tag)
			if fieldName != "" && (allowEmpty || !fieldVal.IsZero()) {
				cond[fieldName] = fieldVal.Interface()
			}
		}
		return nil
	default:
		return errors.ErrUnknownKind
	}
}

func ParseQuery(source interface{}, allowEmpty bool) (string, error) {
	val := GetValue(source)

	switch val.Kind() {
	case reflect.Struct:
		return ParseQueryStruct(val, allowEmpty)
	case reflect.Map:
		return ParseQueryMap(val)
	default:
		return "", errors.ErrUnknownKind
	}

}

func ParseQueryStruct(val reflect.Value, allowEmpty bool) (string, error) {
	if val.Kind() != reflect.Struct {
		return "", errors.ErrUnknownKind
	}

	var query strings.Builder

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldName := val.Type().Field(i).Tag.Get(Tag)
		if fieldName != "" && (allowEmpty || !fieldVal.IsZero()) {
			if i != 0 {
				query.WriteString(" AND ")
			}

			if _, err := query.WriteString(fmt.Sprintf("%s = %v", fieldName, fieldVal.Interface())); err != nil {
				return "", err
			}
		}
	}
	return query.String(), nil
}

func ParseQueryMap(val reflect.Value) (string, error) {
	if val.Kind() != reflect.Map {
		return "", errors.ErrUnknownKind
	}

	var query strings.Builder

	for i, keyVal := range val.MapKeys() {
		valVal := val.MapIndex(keyVal)

		if i != 0 {
			query.WriteString(" AND ")
		}

		condStr := fmt.Sprintf("%s = %v", keyVal.String(), valVal.Interface())
		if _, err := query.WriteString(condStr); err != nil {
			return "", err
		}
	}
	return query.String(), nil
}

func ParseFieldsAndID(source interface{}, allowEmpty bool, named bool) (fields []string, values []string, data []interface{}, id map[string]interface{}, err error) {
	val := GetValue(source)
	id = make(map[string]interface{})

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			fieldName := val.Type().Field(i).Tag.Get(Tag)
			if fieldName != "" && (allowEmpty || !fieldVal.IsZero()) {
				fields = append(fields, fieldName)
				data = append(data, fieldVal.Interface())

				if len(id) == 0 {
					if fieldName == "id" {
						id[fieldName] = fieldVal.Interface()
					} else if fieldName == "uid" || fieldName == "uuid" {
						id[fieldName] = fieldVal.Interface()
					}
				}
				if named {
					values = append(values, ":"+fieldName)
				} else {
					values = append(values, "?")
				}
			}
		}
		return
	default:
		err = errors.ErrUnknownKind
		return
	}
}

func GetValue(source interface{}) reflect.Value {
	val := reflect.ValueOf(source)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}
