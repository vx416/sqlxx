package sqlbuilder

import (
	"github.com/vicxu416/sqlxx/sqlbuilder/query"
)

func (builder *Builder) BuildSelect(table string, query query.Querier) (string, []interface{}, error) {
	whereStr, args, err := query.WhereQueryArgs()
	if err != nil {
		return "", []interface{}{}, err
	}
	selectStr, _, err := query.SelectQueriesArgs()
	if err != nil {
		return "", []interface{}{}, err
	}

	builder.stmt = SelectStmt{
		table:     table,
		whereStr:  whereStr,
		selectStr: selectStr["select"],
		joinStr:   selectStr["join"],
		options:   selectStr["options"],
	}

	return builder.Sql(), args, nil
}

type SelectStmt struct {
	whereStr  string
	selectStr string
	joinStr   string
	options   string
	table     string
}

func (stmt SelectStmt) Keyword(driver driverType) string {
	if stmt.selectStr == "" {
		return "SELECT * FROM"
	}

	return "SELECT " + stmt.selectStr + " FROM"
}

func (stmt SelectStmt) Table(driver driverType) string {
	if stmt.joinStr == "" {
		return stmt.table
	}

	return stmt.table + " " + stmt.joinStr
}

func (stmt SelectStmt) Values(driver driverType) string {
	return ""
}

func (stmt SelectStmt) Where(driver driverType) string {
	return stmt.whereStr
}

func (stmt SelectStmt) End(driver driverType) string {
	return stmt.options
}
