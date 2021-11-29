package builder

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	SkipZero Option = OptionFunc(skipZero)
	Require  Option = OptionFunc(requireValue)
)

type Option interface {
	Check(s string, arg interface{}) (string, bool, error)
}

type OptionFunc func(s string, arg interface{}) (string, bool, error)

func (f OptionFunc) Check(s string, arg interface{}) (string, bool, error) {
	return f(s, arg)
}

func Prefix(prefix string) Option {
	prefix = strings.TrimRight(prefix, ".")
	return PrefixOpt{
		prefix: prefix,
	}
}

type PrefixOpt struct {
	prefix string
}

func (f PrefixOpt) Check(s string, arg interface{}) (string, bool, error) {
	return f.prefix + "." + s, true, nil
}

func skipZero(s string, arg interface{}) (string, bool, error) {
	isZero, err := IsZero(reflect.ValueOf(arg))
	if err != nil {
		return s, false, err
	}

	return s, !isZero, nil
}

func requireValue(s string, arg interface{}) (string, bool, error) {
	isZero, err := IsZero(reflect.ValueOf(arg))
	if err != nil {
		return s, false, err
	}
	if isZero {
		return s, false, errors.New("cannot be empty")
	}
	return s, true, nil
}

func getInArgs(arg interface{}) ([]interface{}, error) {
	val := reflect.ValueOf(arg)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil, errors.New("args is not slice or array")
	}
	res := make([]interface{}, 0, 10)
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if item.Kind() == reflect.Slice || item.Kind() == reflect.Array {
			itemRes, err := getInArgs(item.Interface())
			if err != nil {
				return nil, err
			}
			res = append(res, itemRes...)
			continue
		}
		itemVal := item.Interface()
		if item.Kind() >= reflect.Int8 && item.Kind() <= reflect.Int64 {
			itemVal = int(item.Int())
		}
		if item.Kind() >= reflect.Uint && item.Kind() <= reflect.Uint64 {
			itemVal = int(item.Uint())
		}
		res = append(res, itemVal)
	}

	return res, nil
}

const (
	ReadLock MySQLLock = iota + 1
	WriteLock
)

type MySQLLock uint8

func (value MySQLLock) SQL() string {
	switch value {
	case ReadLock:
		return "LOCK IN SHARE MODE"
	case WriteLock:
		return "FOR UPDATE"
	default:
		return ""
	}
}

type Fields []string

func (value Fields) SQL() string {
	if len(value) == 0 {
		return "*"
	}

	return JoinFields(value, ",")
}

func JoinFields(fields []string, sep string) string {
	str := ""
	for i, field := range fields {
		if i < len(fields)-1 {
			str += field + sep
		} else {
			str += field
		}
	}
	return str
}

// NeedJoin 檢查targetTableName是否含在fields裡
func NeedJoin(fields []string, targetTableName string) bool {
	searchMap := make(map[string]bool)
	for _, field := range fields {
		fieldArray := strings.Split(field, ".")
		tableName := fieldArray[0]
		searchMap[tableName] = true
	}
	return searchMap[targetTableName]
}

type Pagination struct {
	PerPage int
	Page    int
}

func (q Pagination) LimitOffset() (int, int) {
	if q.PerPage <= 0 {
		return 0, 0
	}
	if q.Page <= 0 {
		return q.PerPage, 0
	}

	offset := (q.Page - 1) * q.PerPage
	return q.PerPage, offset
}

func (q Pagination) TotalPages(total int) int {
	totalPages := total / q.PerPage
	if total%q.PerPage != 0 {
		totalPages += 1
	}
	return totalPages
}

func GetElem(data interface{}) reflect.Value {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		return val.Elem()
	}

	return val
}

func IsZero(v reflect.Value) (bool, error) {
	if valid, ok := IsNullable(v); ok {
		return !valid, nil
	}

	if t, ok := v.Interface().(time.Time); ok {
		return t.IsZero() || t.Unix() == 0, nil
	}

	if a, ok := v.Interface().(driver.Valuer); ok {
		var err error
		arg, err := a.Value()
		if err != nil {
			return false, fmt.Errorf("values get value failed, err:%+v", err)
		}
		v = reflect.ValueOf(arg)
		v = ConvertStringNumber(v)
	}

	switch v.Kind() {
	case reflect.Bool:
		return false, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0, nil
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0, nil
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0, nil
	case reflect.Array:
		return v.Len() == 0, nil
	case reflect.Slice:
		return v.IsNil() || v.Len() == 0, nil
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		return v.IsNil(), nil
	case reflect.String:
		return v.Len() == 0, nil
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				return false, nil
			}
		}
		return true, nil
	default:
		// This should never happens, but will act as a safeguard for
		// later, as a default value doesn't makes sense here.
		return true, nil
	}
}

func IsNullable(v reflect.Value) (isValid, ok bool) {
	typ := v.Type()
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false, false
	}

	for i := 0; i < v.NumField(); i++ {
		fieldTyp := typ.Field(i)
		fieldVal := v.Field(i)
		if fieldVal.Kind() == reflect.Struct && fieldTyp.Anonymous {
			if valid, ok := IsNullable(fieldVal); ok {
				return valid, ok
			}
		}
		if fieldTyp.Name == "Valid" && fieldTyp.Type.Kind() == reflect.Bool {
			return fieldVal.Bool(), true
		}
	}

	return false, false
}

func ConvertStringNumber(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.String {
		str := v.String()
		intValue, err := strconv.Atoi(str)
		if err == nil {
			return reflect.ValueOf(intValue)
		}
	}

	return v
}

func joinFields(fields []string, sep string) string {
	str := ""
	for i, field := range fields {
		if i < len(fields)-1 {
			str += field + sep
		} else {
			str += field
		}
	}
	return str
}
