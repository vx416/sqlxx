package builder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

func NewWhereStmt() *WhereStmt {
	return &WhereStmt{
		args: make([]interface{}, 0, 10),
	}
}

type WhereStmt struct {
	strings.Builder
	args []interface{}
}

func (stmt *WhereStmt) clone() *WhereStmt {
	copyArgs := make([]interface{}, len(stmt.args))
	copy(copyArgs, stmt.args)
	return &WhereStmt{
		Builder: strings.Builder{},
		args:    copyArgs,
	}
}

func (stmt *WhereStmt) Build(s *strings.Builder, args *[]interface{}) error {
	whereStmt := stmt.Builder.String()
	if whereStmt == "" {
		return nil
	}

	if s.Len() > 0 {
		_, err := s.WriteString(" ")
		if err != nil {
			return err
		}
	}
	_, err := s.WriteString("WHERE " + whereStmt)
	if err != nil {
		return err
	}
	*args = append(*args, stmt.args...)
	return nil
}

func (builder *WhereStmt) whereIn(query string, arg interface{}) (string, []interface{}, error) {
	args, err := getInArgs(arg)
	if err != nil {
		return "", nil, err
	}

	query, args, err = sqlx.In(query, args)
	if err != nil {
		return "", nil, err
	}
	return query, args, nil
}

func (builder *WhereStmt) buildWhereStmt(operator string, query string, arg interface{}, in bool, options ...Option) ([]interface{}, error) {
	needAppend := true
	subQuery, ok := arg.(Builder)
	if ok {
		subQueryStr, args, err := subQuery.Build()
		if err != nil {
			return nil, err
		}
		query = strings.Replace(query, "?", subQueryStr, 1)
		err = builder.appendWhere(operator, query, args...)
		if err != nil {
			return nil, err
		}
		return args, nil
	}

	for _, option := range options {
		var (
			err    error
			append bool
		)
		query, append, err = option.Check(query, arg)
		if err != nil {
			return nil, err
		}
		needAppend = needAppend && append
	}

	if needAppend {
		if in {
			query, args, err := builder.whereIn(query, arg)
			if err != nil {
				return nil, err
			}
			err = builder.appendWhere(operator, query, args...)
			if err != nil {
				return nil, err
			}
			return args, nil
		}
		err := builder.appendWhere(operator, query, arg)
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, err
		}

		return []interface{}{arg}, nil
	}
	return nil, nil
}

func (builder *WhereStmt) appendWhere(operator string, query string, args ...interface{}) error {
	var err error
	if builder.Len() == 0 {
		_, err = builder.WriteString(query)
	} else {
		_, err = builder.WriteString(fmt.Sprintf(" %s %s", operator, query))
	}
	if err == nil && len(args) > 0 {
		if args[0] == nil {
			return nil
		}
		builder.args = append(builder.args, args...)
	}

	return err
}

func (builder *WhereStmt) appendStruct(s interface{}, options ...Option) ([]interface{}, error) {
	val := GetElem(s)
	if val.Kind() != reflect.Struct {
		return nil, errors.New("input should be struct")
	}
	res := make([]interface{}, 0, 10)
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldT := valType.Field(i)
		sqlTag := fieldT.Tag.Get("sql")
		if sqlTag == "" {
			continue
		}
		qp, err := newQueryParams("AND", sqlTag, val.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		args, err := qp.append(builder, options...)
		if err != nil {
			return nil, err
		}
		res = append(res, args...)
	}

	return res, nil
}

func newQueryParams(cond, tagString string, val interface{}) (*queryParam, error) {
	tagParams := strings.Split(tagString, ";")
	if len(tagParams) == 0 {
		return nil, fmt.Errorf("sql tag(%s) invalid", tagString)
	}
	queryDetails := make(map[string]string)

	for _, tagParam := range tagParams {
		tagParam = strings.TrimSpace(tagParam)
		kv := strings.Split(tagParam, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("sql tag(%s) invalid", tagParam)
		}
		queryDetails[kv[0]] = kv[1]
	}

	if queryDetails["col"] == "" {
		return nil, fmt.Errorf("sql tag(%s) invalid, col cannot be empty", tagString)
	}

	op := strings.ToUpper(queryDetails["op"])
	if op == "" {
		op = "="
	}
	in := false
	if strings.EqualFold(op, "IN") {
		in = true
		op = "IN"
	}
	if strings.EqualFold(op, "NOTIN") {
		in = true
		op = "NOT IN"
	}
	if strings.Contains(op, "%") {
		s, ok := val.(string)
		if !ok {
			return nil, errors.New("like operator should used on string type")
		}
		if s != "" {
			val = strings.Replace(op, "{}", s, 1)
		}
		op = "LIKE"
	}

	qp := &queryParam{
		cond: cond,
		arg:  val,
		col:  queryDetails["col"],
		op:   op,
		in:   in,
	}

	return qp, nil
}

type queryParam struct {
	col  string
	op   string
	in   bool
	cond string
	arg  interface{}
}

func (qp *queryParam) append(where *WhereStmt, options ...Option) ([]interface{}, error) {
	if qp.in {
		return where.buildWhereStmt(qp.cond, fmt.Sprintf("%s %s (?)", qp.col, qp.op), qp.arg, qp.in, options...)
	}

	return where.buildWhereStmt(qp.cond, fmt.Sprintf("%s %s ?", qp.col, qp.op), qp.arg, qp.in, options...)
}

func NewOtherStmt() *OtherStmt {
	return &OtherStmt{
		group: make([]string, 0, 10),
		order: make([]string, 0, 10),
	}
}

type OtherStmt struct {
	group  []string
	order  []string
	limit  int
	offset int
	lock   string
}

func (stmt *OtherStmt) clone() *OtherStmt {
	copyGroup := make([]string, len(stmt.group))
	copyOrder := make([]string, len(stmt.order))
	copy(copyGroup, stmt.group)
	copy(copyOrder, stmt.order)

	return &OtherStmt{
		group:  copyGroup,
		order:  copyOrder,
		limit:  stmt.limit,
		offset: stmt.offset,
		lock:   stmt.lock,
	}
}

func (stmt *OtherStmt) Build(s *strings.Builder, args *[]interface{}) error {
	groupBy := joinFields(stmt.group, ", ")
	orderBy := joinFields(stmt.order, ", ")
	if s.Len() > 0 && (groupBy != "" || orderBy != "" || stmt.limit != 0 || stmt.offset != 0 || stmt.lock != "") {
		_, err := s.WriteString(" ")
		if err != nil {
			return err
		}
	}
	space := ""
	if groupBy != "" {
		_, err := s.WriteString("GROUP BY " + groupBy)
		if err != nil {
			return err
		}
		space = " "
	}
	if orderBy != "" {
		_, err := s.WriteString(space + "ORDER BY " + orderBy)
		if err != nil {
			return err
		}
		space = " "
	}
	if stmt.limit > 0 {
		_, err := s.WriteString(space + fmt.Sprintf("LIMIT %d", stmt.limit))
		if err != nil {
			return err
		}
		space = " "
	}
	if stmt.offset > 0 {
		_, err := s.WriteString(space + fmt.Sprintf("OFFSET %d", stmt.offset))
		if err != nil {
			return err
		}
		space = " "
	}
	if stmt.lock != "" {
		_, err := s.WriteString(space + string(stmt.lock))
		if err != nil {
			return err
		}
	}

	return nil
}

func (stmt *OtherStmt) groupBy(ss ...string) error {
	stmt.group = append(stmt.group, ss...)
	return nil
}

func (stmt *OtherStmt) orderBy(ss ...string) error {
	stmt.order = append(stmt.order, ss...)
	return nil
}

func (stmt *OtherStmt) limitOffset(limit, offset int) error {
	stmt.limit = limit
	stmt.offset = offset
	return nil
}
