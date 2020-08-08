package parsers

import (
	"reflect"

	"github.com/vicxu416/sqlxx/errors"
)

var Tag = "db"

func New(data interface{}, allEmpty bool) (*Parser, error) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	parser := &Parser{
		source:      val,
		Fields:      make([]string, 0, 1),
		Values:      make([][]string, 0, 1),
		NamedValues: make([][]string, 0, 1),
		Data:        make([]interface{}, 0, 1),
		kind:        val.Kind(),
	}

	if err := parser.Parse(allEmpty); err != nil {
		return nil, err
	}

	return parser, nil
}

type Parser struct {
	source      reflect.Value
	Fields      []string
	Values      [][]string
	NamedValues [][]string
	Data        []interface{}
	kind        reflect.Kind
	hasCache    bool
}

func (parser *Parser) cache(fields []string, values [][]string, data []interface{}, named bool) {
	parser.Fields = fields
	if named {
		parser.NamedValues = values
	} else {
		parser.Values = values
	}

	parser.Data = data
}

func (parser *Parser) Parse(allowEmpty bool) error {
	if parser.hasCache {
		return nil
	}

	switch parser.kind {
	case reflect.Struct:
		if err := parseStruct(parser.source, &parser.Fields, &parser.Values, &parser.NamedValues, &parser.Data, allowEmpty); err != nil {
			return err
		}
	case reflect.Map:
		if err := parseMap(parser.source, &parser.Fields, &parser.Values, &parser.NamedValues, &parser.Data); err != nil {
			return err
		}
	case reflect.Slice:
		if err := parseSlice(parser.source, &parser.Fields, &parser.Values, &parser.NamedValues, &parser.Data, allowEmpty); err != nil {
			return err
		}
	}
	parser.hasCache = true
	return nil
}

func parseSlice(val reflect.Value, fields *[]string, values *[][]string, namedValues *[][]string, data *[]interface{}, allowEmpty bool) error {
	if val.Kind() != reflect.Slice {
		return errors.ErrUnknownKind
	}

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		switch item.Kind() {
		case reflect.Struct:
			if err := parseStruct(item, fields, values, namedValues, data, allowEmpty); err != nil {
				return err
			}
		case reflect.Map:
			if err := parseMap(item, fields, values, namedValues, data); err != nil {
				return err
			}
		default:
			return errors.ErrUnknownKind
		}
	}

	return nil
}

func parseStruct(val reflect.Value, fields *[]string, values *[][]string, namedValues *[][]string, data *[]interface{}, allowEmpty bool) error {
	if val.Kind() != reflect.Struct {
		return errors.ErrUnknownKind
	}

	subValues := make([]string, 0, val.NumField())
	namedSubValues := make([]string, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldName := val.Type().Field(i).Tag.Get(Tag)
		if fieldName != "" && (allowEmpty || !fieldVal.IsZero()) {
			*fields = append(*fields, fieldName)
			namedSubValues = append(namedSubValues, ":"+fieldName)
			subValues = append(subValues, "?")
			*data = append(*data, fieldVal.Interface())
		}
	}
	*values = append(*values, subValues)
	*namedValues = append(*namedValues, namedSubValues)

	return nil
}

func parseMap(val reflect.Value, fields *[]string, values *[][]string, namedValues *[][]string, data *[]interface{}) error {
	if val.Kind() != reflect.Map {
		return errors.ErrUnknownKind
	}

	subValues := make([]string, 0, len(val.MapKeys()))
	namedSubValues := make([]string, 0, len(val.MapKeys()))

	for _, keyVal := range val.MapKeys() {
		*fields = append(*fields, keyVal.String())
		namedSubValues = append(namedSubValues, ":"+keyVal.String())
		subValues = append(subValues, "?")
		*data = append(*data, val.MapIndex(keyVal).Interface())
	}
	*values = append(*values, subValues)
	*namedValues = append(*namedValues, namedSubValues)

	return nil
}
