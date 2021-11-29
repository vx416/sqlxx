package builder

import (
	"errors"
	"strings"
)

func Delete() *DeleteBuilder {
	return &DeleteBuilder{
		whereStmt: NewWhereStmt(),
	}
}

type DeleteBuilder struct {
	table     string
	whereStmt *WhereStmt
	err       error
}

func (builder *DeleteBuilder) Build() (string, []interface{}, error) {
	if builder.err != nil {
		return "", nil, builder.err
	}
	if builder.table == "" {
		return "", nil, errors.New("table cannot be empty")
	}
	query := &strings.Builder{}
	args := make([]interface{}, 0, 10)

	_, err := query.WriteString("DELETE FROM ")
	if err != nil {
		return "", nil, err
	}
	_, err = query.WriteString(builder.table)
	if err != nil {
		return "", nil, err
	}

	err = builder.whereStmt.Build(query, &args)
	if err != nil {
		return "", nil, err
	}
	return query.String(), args, nil
}

func (builder *DeleteBuilder) Table(t string) *DeleteBuilder {
	builder.table = t
	return builder
}

func (builder *DeleteBuilder) Clone() *DeleteBuilder {
	cloned := &DeleteBuilder{
		whereStmt: builder.whereStmt.clone(),
		err:       builder.err,
	}
	_, err := cloned.whereStmt.WriteString(builder.whereStmt.String())
	builder.setErr(err)
	return cloned
}

func (builder *DeleteBuilder) From(table string) *DeleteBuilder {
	if builder.err != nil {
		return builder
	}
	builder.table = table
	return builder
}

func (builder *DeleteBuilder) And(query string, arg interface{}, options ...Option) *DeleteBuilder {
	builder.appendWhereStmt("AND", query, arg, false, options...)
	return builder
}

func (builder *DeleteBuilder) Or(query string, arg interface{}, options ...Option) *DeleteBuilder {
	builder.appendWhereStmt("OR", query, arg, false, options...)
	return builder
}

func (builder *DeleteBuilder) AndIn(query string, arg interface{}, options ...Option) *DeleteBuilder {
	builder.appendWhereStmt("AND", query, arg, true, options...)
	return builder
}

func (builder *DeleteBuilder) OrIn(query string, arg interface{}, options ...Option) *DeleteBuilder {
	builder.appendWhereStmt("OR", query, arg, true, options...)
	return builder
}

func (builder *DeleteBuilder) Where(st interface{}, options ...Option) *DeleteBuilder {
	if builder.err != nil {
		return builder
	}
	if tabler, ok := st.(Tabler); ok {
		if builder.table == "" {
			builder.table = tabler.TableName()
		}
	}
	_, err := builder.whereStmt.appendStruct(st, options...)
	builder.setErr(err)
	return builder
}

func (builder *DeleteBuilder) appendWhereStmt(op, query string, arg interface{}, in bool, options ...Option) {
	if builder.err != nil {
		return
	}
	_, err := builder.whereStmt.buildWhereStmt(op, query, arg, in, options...)
	builder.setErr(err)
}

func (builder *DeleteBuilder) setErr(err error) {
	if err != nil && builder.err == nil {
		builder.err = err
	}
}
