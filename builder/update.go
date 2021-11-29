package builder

import (
	"strings"
)

func Update() *UpdateBuilder {
	return &UpdateBuilder{
		updateStmt: NewUpdateStmt(),
		whereStmt:  NewWhereStmt(),
	}
}

type UpdateBuilder struct {
	updateStmt *UpdateStmt
	whereStmt  *WhereStmt
	err        error
}

func (builder *UpdateBuilder) Build() (string, []interface{}, error) {
	if builder.err != nil {
		return "", nil, builder.err
	}
	query := &strings.Builder{}
	args := make([]interface{}, 0, 10)

	err := builder.updateStmt.Build(query, &args)
	if err != nil {
		return "", nil, err
	}

	err = builder.whereStmt.Build(query, &args)
	if err != nil {
		return "", nil, err
	}
	return query.String(), args, nil
}

func (builder *UpdateBuilder) Table(t string) *UpdateBuilder {
	builder.updateStmt.table = t
	return builder
}

func (builder *UpdateBuilder) Clone() *UpdateBuilder {
	cloned := &UpdateBuilder{
		updateStmt: builder.updateStmt.clone(),
		whereStmt:  builder.whereStmt.clone(),
		err:        builder.err,
	}
	_, err := cloned.updateStmt.WriteString(builder.updateStmt.String())
	builder.setErr(err)
	_, err = cloned.whereStmt.WriteString(builder.whereStmt.String())
	builder.setErr(err)
	return cloned
}

func (builder *UpdateBuilder) Set(set string, arg interface{}, options ...Option) *UpdateBuilder {
	if builder.err != nil {
		return builder
	}
	err := builder.updateStmt.set(set, arg, options...)

	builder.setErr(err)
	return builder
}

func (builder *UpdateBuilder) SetWith(arg interface{}, opts ...Option) *UpdateBuilder {
	if builder.err != nil {
		return builder
	}
	if tabler, ok := arg.(Tabler); ok {
		if builder.updateStmt.table == "" {
			builder.updateStmt.setTable(tabler.TableName())
		}
	}
	err := builder.updateStmt.updateWith(arg, opts...)
	builder.setErr(err)
	return builder
}

func (builder *UpdateBuilder) And(query string, arg interface{}, options ...Option) *UpdateBuilder {
	builder.appendWhereStmt("AND", query, arg, false, options...)
	return builder
}

func (builder *UpdateBuilder) Or(query string, arg interface{}, options ...Option) *UpdateBuilder {
	builder.appendWhereStmt("OR", query, arg, false, options...)
	return builder
}

func (builder *UpdateBuilder) AndIn(query string, arg interface{}, options ...Option) *UpdateBuilder {
	builder.appendWhereStmt("AND", query, arg, true, options...)
	return builder
}

func (builder *UpdateBuilder) OrIn(query string, arg interface{}, options ...Option) *UpdateBuilder {
	builder.appendWhereStmt("OR", query, arg, true, options...)
	return builder
}

func (builder *UpdateBuilder) Where(st interface{}, options ...Option) *UpdateBuilder {
	if builder.err != nil {
		return builder
	}
	if tabler, ok := st.(Tabler); ok {
		if builder.updateStmt.table == "" {
			builder.updateStmt.setTable(tabler.TableName())
		}
	}
	_, err := builder.whereStmt.appendStruct(st, options...)
	builder.setErr(err)
	return builder
}

func (builder *UpdateBuilder) appendWhereStmt(op, query string, arg interface{}, in bool, options ...Option) {
	if builder.err != nil {
		return
	}
	_, err := builder.whereStmt.buildWhereStmt(op, query, arg, in, options...)
	builder.setErr(err)
}

func (builder *UpdateBuilder) setErr(err error) {
	if err != nil && builder.err == nil {
		builder.err = err
	}
}
