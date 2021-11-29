package builder

import (
	"errors"
	"fmt"
	"strings"
)

type Tabler interface {
	TableName() string
}

type Builder interface {
	Build() (string, []interface{}, error)
}

type unionQuery struct {
	union string
	*QueryBuilder
}

func Query() *QueryBuilder {
	return &QueryBuilder{
		selectStmt: NewSelectStmt(),
		whereStmt:  NewWhereStmt(),
		otherStmt:  NewOtherStmt(),
		unions:     make([]*unionQuery, 0, 10),
	}
}

type QueryBuilder struct {
	selectStmt *SelectStmt
	whereStmt  *WhereStmt
	otherStmt  *OtherStmt
	err        error
	unions     []*unionQuery
}

func (builder *QueryBuilder) Clone() *QueryBuilder {
	copyUnions := make([]*unionQuery, 0, cap(builder.unions))
	for _, union := range builder.unions {
		copyUnions = append(copyUnions, &unionQuery{
			union:        union.union,
			QueryBuilder: union.Clone(),
		})
	}

	cloned := &QueryBuilder{
		selectStmt: builder.selectStmt.clone(),
		whereStmt:  builder.whereStmt.clone(),
		otherStmt:  builder.otherStmt.clone(),
		err:        builder.err,
		unions:     copyUnions,
	}
	_, err := cloned.selectStmt.WriteString(builder.selectStmt.String())
	builder.setErr(err)
	_, err = cloned.selectStmt.joins.WriteString(builder.selectStmt.joins.String())
	builder.setErr(err)
	_, err = cloned.whereStmt.WriteString(builder.whereStmt.String())
	builder.setErr(err)
	return cloned
}

func (builder *QueryBuilder) Build() (string, []interface{}, error) {
	if builder.err != nil {
		return "", nil, builder.err
	}
	query := &strings.Builder{}
	args := make([]interface{}, 0, 10)
	err := builder.selectStmt.Build(query, &args)
	if err != nil {
		return "", nil, err
	}
	err = builder.whereStmt.Build(query, &args)
	if err != nil {
		return "", nil, err
	}

	err = builder.otherStmt.Build(query, &args)
	if err != nil {
		return "", nil, err
	}

	queryS := query.String()
	if len(builder.unions) > 0 {
		queryS = fmt.Sprintf("(%s)", queryS)
	}
	for _, union := range builder.unions {
		unionQuery, unionArgs, err := union.Build()
		if err != nil {
			return "", nil, err
		}
		queryS += fmt.Sprintf(" %s (%s)", union.union, unionQuery)
		args = append(args, unionArgs...)
	}
	return queryS, args, nil
}

func (builder *QueryBuilder) Union(other *QueryBuilder) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	builder.unions = append(builder.unions, &unionQuery{
		union:        "UNION",
		QueryBuilder: other,
	})
	return builder
}

func (builder *QueryBuilder) UnionAll(other *QueryBuilder) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	builder.unions = append(builder.unions, &unionQuery{
		union:        "UNION ALL",
		QueryBuilder: other,
	})
	return builder
}

func (builder *QueryBuilder) Select(ss ...string) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	err := builder.selectStmt.appendSelect(ss...)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) Count() *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	err := builder.selectStmt.count()
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) From(s string, args ...interface{}) *QueryBuilder {
	if builder.err != nil {
		return builder
	}

	var (
		subQuery Builder = nil
		ok       bool
	)

	if len(args) > 0 {
		subQuery, ok = args[0].(Builder)
		if !ok {
			builder.err = errors.New("from args can only be pointer of Query")
			return builder
		}
	}

	err := builder.selectStmt.setFrom(s, subQuery)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) Join(s string, args ...interface{}) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	var subQuery Builder = nil
	if len(args) > 0 {
		if q, ok := args[0].(Builder); ok {
			subQuery = q
		}
	}

	err := builder.selectStmt.join(Join, s, subQuery)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) LeftJoin(s string, args ...interface{}) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	var subQuery Builder = nil
	if len(args) > 0 {
		if q, ok := args[0].(Builder); ok {
			subQuery = q
		}
	}

	err := builder.selectStmt.join(LeftJoin, s, subQuery)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) RightJoin(s string, args ...interface{}) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	var subQuery Builder = nil
	if len(args) > 0 {
		if q, ok := args[0].(Builder); ok {
			subQuery = q
		}
	}

	err := builder.selectStmt.join(RightJoin, s, subQuery)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) OuterJoin(s string, args ...interface{}) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	var subQuery Builder = nil
	if len(args) > 0 {
		if q, ok := args[0].(Builder); ok {
			subQuery = q
		}
	}

	err := builder.selectStmt.join(OuterJoin, s, subQuery)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) And(query string, arg interface{}, options ...Option) *QueryBuilder {
	builder.appendWhereStmt("AND", query, arg, false, options...)
	return builder
}

func (builder *QueryBuilder) Or(query string, arg interface{}, options ...Option) *QueryBuilder {
	builder.appendWhereStmt("OR", query, arg, false, options...)
	return builder
}

func (builder *QueryBuilder) AndIn(query string, arg interface{}, options ...Option) *QueryBuilder {
	builder.appendWhereStmt("AND", query, arg, true, options...)
	return builder
}

func (builder *QueryBuilder) OrIn(query string, arg interface{}, options ...Option) *QueryBuilder {
	builder.appendWhereStmt("OR", query, arg, true, options...)
	return builder
}

func (builder *QueryBuilder) Where(st interface{}, options ...Option) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	if tabler, ok := st.(Tabler); ok {
		if builder.selectStmt.from == "" {
			builder.From(tabler.TableName())
		}
	}
	_, err := builder.whereStmt.appendStruct(st, options...)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) GroupBy(ss ...string) *QueryBuilder {
	if builder.err != nil {
		return builder
	}

	err := builder.otherStmt.groupBy(ss...)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) OrderBy(ss ...string) *QueryBuilder {
	if builder.err != nil {
		return builder
	}

	err := builder.otherStmt.orderBy(ss...)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) LimitOffset(limit, offset int) *QueryBuilder {
	if builder.err != nil {
		return builder
	}

	err := builder.otherStmt.limitOffset(limit, offset)
	builder.setErr(err)
	return builder
}

func (builder *QueryBuilder) Lock(s string) *QueryBuilder {
	if builder.err != nil {
		return builder
	}
	builder.otherStmt.lock = s
	return builder
}

func (builder *QueryBuilder) appendWhereStmt(op, query string, arg interface{}, in bool, options ...Option) {
	if builder.err != nil {
		return
	}
	_, err := builder.whereStmt.buildWhereStmt(op, query, arg, in, options...)
	builder.setErr(err)
}

func (builder *QueryBuilder) setErr(err error) {
	if err != nil && builder.err == nil {
		builder.err = err
	}
}
