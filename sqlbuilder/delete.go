package sqlbuilder

import (
	"strings"

	"github.com/vicxu416/sqlxx/sqlbuilder/query"
)

func (builder *Builder) BuildDelete(table string, usingTable []string, query query.Querier) (string, []interface{}, error) {
	whereStr, args, err := query.WhereQueryArgs()
	if err != nil {
		return "", []interface{}{}, err
	}

	builder.stmt = DeleteStmt{
		whereStr:   whereStr,
		table:      table,
		usingTable: usingTable,
	}
	return builder.Sql(), args, nil
}

type DeleteStmt struct {
	whereStr   string
	table      string
	usingTable []string
}

func (stmt DeleteStmt) Keyword(driver driverType) string {
	return "DELETE FROM"
}

func (stmt DeleteStmt) Table(driver driverType) string {
	if len(stmt.usingTable) == 0 {
		return stmt.table
	}
	return stmt.table + " USING " + strings.Join(stmt.usingTable, ",")
}

func (stmt DeleteStmt) Values(driver driverType) string {
	return ""
}

func (stmt DeleteStmt) Where(driver driverType) string {
	return stmt.whereStr
}

func (stmt DeleteStmt) End(driver driverType) string {
	return ""
}
