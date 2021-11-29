package builder

import (
	"strings"
)

const (
	MSShareLock = "LOCK IN SHARE MODE"
	MSWRITELOCK = "FOR UPDATE"
)

type JoinType string

const (
	LeftJoin  JoinType = "LEFT JOIN"
	RightJoin JoinType = "RIGHT JOIN"
	OuterJoin JoinType = "OUTTER JOIN"
	Join      JoinType = "JOIN"
)

func NewSelectStmt() *SelectStmt {
	return &SelectStmt{
		args: make([]interface{}, 0, 10),
	}
}

type SelectStmt struct {
	strings.Builder
	from  string
	joins strings.Builder
	args  []interface{}
}

func (stmt *SelectStmt) clone() *SelectStmt {
	copyArgs := make([]interface{}, len(stmt.args))
	copy(copyArgs, stmt.args)
	return &SelectStmt{
		Builder: strings.Builder{},
		from:    stmt.from,
		joins:   strings.Builder{},
		args:    copyArgs,
	}
}

func (stmt *SelectStmt) Build(s *strings.Builder, args *[]interface{}) error {
	selectStmt := stmt.Builder.String()
	if selectStmt == "" {
		selectStmt = "*"
	}

	if s.Len() > 0 {
		_, err := s.WriteString(" ")
		if err != nil {
			return err
		}
	}

	_, err := s.WriteString("SELECT " + selectStmt + " FROM " + stmt.from)
	if err != nil {
		return err
	}
	joinStmt := stmt.joins.String()
	if joinStmt != "" {
		s.WriteString(" " + joinStmt)
	}
	*args = append(*args, stmt.args...)

	return nil
}

func (stmt *SelectStmt) setFrom(from string, subQuery Builder) error {
	stmt.from = from

	if subQuery != nil {
		subQueryStr, args, err := subQuery.Build()
		if err != nil {
			return err
		}
		stmt.from = strings.Replace(stmt.from, "?", subQueryStr, 1)
		stmt.args = append(stmt.args, args...)
	}

	return nil
}

func (stmt *SelectStmt) appendSelect(ss ...string) error {
	s := joinFields(ss, ", ")

	if stmt.Len() > 0 {
		_, err := stmt.WriteString(", " + s)
		return err
	}
	_, err := stmt.WriteString(s)
	return err
}

func (stmt *SelectStmt) join(joinType JoinType, s string, subQuery Builder) error {
	if subQuery != nil {
		subQS, subArgs, err := subQuery.Build()
		if err != nil {
			return err
		}
		s = strings.Replace(s, "?", subQS, 1)
		stmt.args = append(stmt.args, subArgs...)
	}

	if stmt.joins.Len() > 0 {
		_, err := stmt.joins.WriteString(" " + string(joinType) + " " + s)
		return err
	}
	_, err := stmt.joins.WriteString(string(joinType) + " " + s)
	return err
}

func (stmt *SelectStmt) count() error {
	stmt.Reset()
	_, err := stmt.WriteString("COUNT(1)")
	return err
}
