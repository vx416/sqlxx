package builder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Insert() *InsertBuilder {
	return &InsertBuilder{
		rows: make([]map[string]interface{}, 0, 10),
	}
}

type InsertBuilder struct {
	table string
	rows  []map[string]interface{}
	err   error
}

func (builder *InsertBuilder) Build() (string, []interface{}, error) {
	if builder.err != nil {
		return "", nil, builder.err
	}
	if len(builder.rows) == 0 {
		return "", nil, errors.New("rows is empty")
	}

	fields := make([]string, 0, 10)
	args := make([]interface{}, 0, 10)
	valuesStr := strings.Builder{}

	for _, row := range builder.rows {
		if len(fields) == 0 {
			for k, v := range row {
				fields = append(fields, k)
				args = append(args, v)
			}
		} else {
			if len(row) != len(fields) {
				return "", nil, errors.New("rows size is different")
			}
			for _, field := range fields {
				args = append(args, row[field])
			}
		}

		if valuesStr.Len() > 0 {
			valuesStr.WriteString(", ")
		}
		values := strings.TrimRight(strings.Repeat("?, ", len(fields)), ", ")
		_, err := valuesStr.WriteString("(")
		if err != nil {
			return "", nil, err
		}
		_, err = valuesStr.WriteString(values)
		if err != nil {
			return "", nil, err
		}
		_, err = valuesStr.WriteString(")")
		if err != nil {
			return "", nil, err
		}
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", builder.table, joinFields(fields, ", "), valuesStr.String()), args, nil
}

func (builder *InsertBuilder) Clone() *InsertBuilder {
	newRows := make([]map[string]interface{}, len(builder.rows))

	for i, row := range builder.rows {
		newRow := make(map[string]interface{})
		for k, v := range row {
			newRow[k] = v
		}
		newRows[i] = newRow

	}

	cloned := &InsertBuilder{
		table: builder.table,
		rows:  newRows,
		err:   builder.err,
	}
	return cloned
}

func (builder *InsertBuilder) InsertRows(data interface{}, options ...Option) *InsertBuilder {
	elem := GetElem(data)
	switch elem.Kind() {
	case reflect.Map:
		builder.rows = append(builder.rows, elem.Interface().(map[string]interface{}))
	case reflect.Struct:
		if tabler, ok := data.(Tabler); ok {
			builder.table = tabler.TableName()
		}

		dataMap, err := builder.structToMap(elem, options...)
		if err != nil {
			builder.setErr(err)
			return builder
		}
		builder.rows = append(builder.rows, dataMap)
	case reflect.Slice:
		builder.batchInsert(elem, options...)
	default:
		builder.setErr(errors.New("row kind is not struct or map"))
	}

	return builder
}

func (builder *InsertBuilder) Table(t string) *InsertBuilder {
	builder.table = t
	return builder
}

func (builder *InsertBuilder) batchInsert(val reflect.Value, options ...Option) {
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		builder.InsertRows(item.Interface(), options...)
		if builder.err != nil {
			return
		}
	}
}

func (builder *InsertBuilder) structToMap(val reflect.Value, options ...Option) (map[string]interface{}, error) {
	dataType := val.Type()
	fieldsLen := val.NumField()
	dataMap := make(map[string]interface{})
	var err error
	for i := 0; i < fieldsLen; i++ {
		fieldVal := val.Field(i)
		dbCol := dataType.Field(i).Tag.Get("db")
		if dbCol != "" {
			ok := true
			for _, opt := range options {
				dbCol, ok, err = opt.Check(dbCol, fieldVal.Interface())
				if err != nil {
					return nil, err
				}
			}
			if ok {
				zero, err := IsZero(fieldVal)
				if err != nil {
					return nil, err
				}
				if !zero {
					dataMap[dbCol] = fieldVal.Interface()
				}
			}
		}

	}

	return dataMap, nil
}

func (builder *InsertBuilder) setErr(err error) {
	if err != nil && builder.err == nil {
		builder.err = err
	}
}
