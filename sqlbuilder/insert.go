package sqlbuilder

import (
	"strings"

	"github.com/vicxu416/sqlxx/sqlbuilder/parsers"
)

func (builder *Builder) BuildBulkInsert(table string, source interface{}) (string, []interface{}, error) {
	parser, err := parsers.New(source, false)
	if err != nil {
		return "nil", []interface{}{}, err
	}
	builder.stmt = InsertStmt{
		fields: parser.Fields,
		values: parser.Values,
		table:  table,
	}

	if err != nil {
		return "", []interface{}{}, err
	}
	return builder.Sql(), parser.Data, nil
}

func (builder *Builder) BuildInsert(table string, source interface{}) (string, []interface{}, error) {
	parser, err := parsers.New(source, false)
	if err != nil {
		return "", []interface{}{}, err
	}
	builder.stmt = InsertStmt{
		fields: parser.Fields,
		values: parser.Values,
		table:  table,
	}
	return builder.Sql(), parser.Data, nil
}

type InsertStmt struct {
	fields []string
	values [][]string
	table  string
}

func (stmt InsertStmt) Keyword(driver driverType) string {
	return "INSERT INTO"
}

func (stmt InsertStmt) Table(driver driverType) string {
	return stmt.table
}

func (stmt InsertStmt) Values(driver driverType) string {
	var builder strings.Builder

	_ = builder.WriteByte('(')
	_, _ = builder.WriteString(strings.Join(stmt.fields, ", "))
	_ = builder.WriteByte(')')
	_, _ = builder.WriteString(" VALUES ")

	for i, value := range stmt.values {
		_ = builder.WriteByte('(')
		_, _ = builder.WriteString(strings.Join(value, ", "))
		_ = builder.WriteByte(')')
		if i != len(stmt.values)-1 {
			_ = builder.WriteByte(',')
		}
	}

	return builder.String()
}

func (stmt InsertStmt) Where(driver driverType) string {
	return ""
}

func (stmt InsertStmt) End(driver driverType) string {
	if driver == pg {
		return "returning id"
	}
	return ""
}
