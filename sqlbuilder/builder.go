package sqlbuilder

import (
	"bytes"
)

type driverType string

const (
	pg driverType = "postgres"
)

type Statement interface {
	Keyword(driverType) string
	Table(driverType) string
	Values(driverType) string
	Where(driverType) string
	End(driverType) string
}

func New(dbType string) *Builder {
	return &Builder{
		dbType: driverType(dbType),
		Buffer: bytes.Buffer{},
	}
}

type Builder struct {
	stmt Statement
	bytes.Buffer
	dbType driverType
}

func (builder Builder) Sql() string {
	builder.Reset()
	builder.WriteString(builder.stmt.Keyword(builder.dbType))
	builder.WriteRune(' ')
	builder.WriteString(builder.stmt.Table(builder.dbType))

	valuseStmt := builder.stmt.Values(builder.dbType)
	if valuseStmt != "" {
		builder.WriteRune(' ')
		builder.WriteString(valuseStmt)
	}

	whereStmt := builder.stmt.Where(builder.dbType)
	if whereStmt != "" {
		builder.WriteRune(' ')
		builder.WriteString("WHERE ")
		builder.WriteString(whereStmt)
	}

	endStmt := builder.stmt.End(builder.dbType)
	if endStmt != "" {
		builder.WriteRune(' ')
		builder.WriteString(endStmt)
	}
	builder.WriteRune(';')

	return builder.String()
}
