package sqlbuilder

import (
	"fmt"
	"strings"

	"github.com/vicxu416/sqlxx/sqlbuilder/query"

	"github.com/vicxu416/sqlxx/sqlbuilder/parsers"
)

func (builder *Builder) BuildUpdateStruct(table string, source interface{}) (string, []interface{}, error) {
	fields, _, data, id, err := parsers.ParseFieldsAndID(source, false, false)
	if err != nil {
		return "", []interface{}{}, err
	}
	var whereStr string

	if len(id) > 0 {
		for k, v := range id {
			whereStr = fmt.Sprintf("%s = %+v", k, v)
			break
		}
	}

	builder.stmt = UpdateStmt{
		fields:   fields,
		whereStr: whereStr,
		table:    table,
	}
	return builder.Sql(), data, nil
}

func (builder *Builder) BuildUpdate(table string, source interface{}, query query.Querier) (string, []interface{}, error) {
	parser, err := parsers.New(source, false)
	if err != nil {
		return "", []interface{}{}, err
	}

	args := parser.Data

	whereStr, whereArgs, err := query.WhereQueryArgs()
	if err != nil {
		return "", []interface{}{}, err
	}

	args = append(args, whereArgs...)

	builder.stmt = UpdateStmt{
		fields:   parser.Fields,
		whereStr: whereStr,
		table:    table,
	}
	return builder.Sql(), args, nil
}

type UpdateStmt struct {
	fields   []string
	whereStr string
	table    string
}

func (stmt UpdateStmt) Keyword(driver driverType) string {
	return "UPDATE"
}

func (stmt UpdateStmt) Table(driver driverType) string {
	return stmt.table
}

func (stmt UpdateStmt) Values(driver driverType) string {
	builder := strings.Builder{}

	for i, field := range stmt.fields {
		if i == 0 {
			builder.WriteString("SET ")
		}
		if i != 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(field + " = " + "?")
	}

	return builder.String()
}

func (stmt UpdateStmt) Where(driver driverType) string {
	return stmt.whereStr
}

func (stmt UpdateStmt) End(driver driverType) string {
	return ""
}
