package sqlbuilder

import (
	"strings"
)

type driverType string

const (
	pg driverType = "postgres"
)

type Statement interface {
	Keyword(driverType) string
	Table(driverType) string
	Values(driverType) string
	Where(driverType) Where
	End(driverType) string
}

type Where interface {
}

type Builder struct {
	stmt Statement
	strings.Builder
}

func (builder Builder) Sql(dbtype string) string {
	driverType := driverType(dbtype)
	builder.Reset()
	builder.WriteString(builder.stmt.Keyword(driverType))
	builder.WriteRune(' ')
	builder.WriteString(builder.stmt.Table(driverType))
	builder.WriteRune(' ')
	builder.WriteString(builder.stmt.Values(driverType))
	builder.WriteRune(' ')

	if builder.stmt.Where(driverType) != nil {

	}
	builder.WriteString(builder.stmt.End(driverType))

	return builder.String()
}
