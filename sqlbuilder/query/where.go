package query

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vicxu416/sqlxx/sqlbuilder/parsers"
)

func NewWhere() *Where {
	return &Where{
		Builder: strings.Builder{},
		args:    make([]interface{}, 0, 10),
	}
}

type Where struct {
	strings.Builder
	args    []interface{}
	lastErr error
}

func (where *Where) QueryArgs() (string, []interface{}, error) {
	return where.String(), where.args, where.lastErr
}

func (where *Where) Where(query string, args ...interface{}) *Where {
	return where.Write(query, "AND", args...)
}

func (where *Where) Or(query string, args ...interface{}) *Where {
	return where.Write(query, "OR", args...)
}

func (where *Where) AndStruct(arg interface{}) *Where {
	query, err := parsers.ParseQuery(arg, false)
	if err != nil {
		where.lastErr = err
	}

	return where.Write(query, "AND")
}

func (where *Where) OrStruct(arg interface{}) *Where {
	query, err := parsers.ParseQuery(arg, false)
	if err != nil {
		where.lastErr = err
	}

	return where.Write(query, "OR")
}

func (where *Where) In(column string, arg interface{}) *Where {
	query, args, err := sqlx.In(column+" IN (?)", arg)
	if err != nil {
		where.lastErr = err
	}
	return where.Where(query, args...)
}

func (where *Where) Write(query, andOR string, args ...interface{}) *Where {
	if where.Len() != 0 {
		where.WriteString(" " + andOR + " ")
	}
	where.WriteRune('(')
	where.WriteString(query)
	where.WriteRune(')')

	if len(args) > 0 {
		where.args = append(where.args, args...)
	}
	return where
}
