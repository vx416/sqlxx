package sqlbuilder

import (
	"strings"

	"github.com/vicxu416/sqlxx/sqlbuilder/parsers"
)

func BulkInsert(table string, source interface{}) (*Builder, []interface{}, error) {
	parser, err := parsers.New(source, false)
	if err != nil {
		return nil, nil, err
	}
	insert := InsertStmt{
		fields: parser.Fields,
		values: parser.Values,
		table:  table,
	}

	if err != nil {
		return nil, nil, err
	}
	return &Builder{stmt: insert}, parser.Data, nil

}

func Insert(table string, source interface{}) (*Builder, error) {
	parser, err := parsers.New(source, false)
	if err != nil {
		return nil, err
	}
	insert := InsertStmt{
		fields: parser.Fields,
		values: parser.NamedValues,
		table:  table,
	}
	return &Builder{stmt: insert}, nil
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

func (stmt InsertStmt) Where(driver driverType) Where {
	return nil
}

func (stmt InsertStmt) End(driver driverType) string {
	if driver == pg {
		return "returning id"
	}
	return ""
}
