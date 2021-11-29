package builder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func NewUpdateStmt() *UpdateStmt {
	return &UpdateStmt{
		args: make([]interface{}, 0, 10),
	}
}

type UpdateStmt struct {
	strings.Builder
	table string
	args  []interface{}
}

func (stmt *UpdateStmt) Build(s *strings.Builder, args *[]interface{}) error {
	if stmt.table == "" {
		return errors.New("table cannot be empty")
	}

	_, err := s.WriteString("UPDATE " + stmt.table + " SET ")
	if err != nil {
		return err
	}
	_, err = s.WriteString(stmt.String())
	if err != nil {
		return err
	}
	*args = append(*args, stmt.args...)
	return nil
}

func (stmt *UpdateStmt) set(set string, arg interface{}, options ...Option) error {
	var (
		err error
		ok  bool
	)
	for _, opt := range options {
		set, ok, err = opt.Check(set, arg)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}
	if stmt.Len() > 0 {
		stmt.WriteString(", ")
	}
	_, err = stmt.WriteString(set)
	if err != nil {
		return err
	}
	if arg != nil {
		stmt.args = append(stmt.args, arg)
	}
	return nil
}

func (stmt *UpdateStmt) setTable(table string) {
	stmt.table = table
}

func (stmt *UpdateStmt) updateWith(data interface{}, options ...Option) error {
	if dataMap, ok := data.(map[string]interface{}); ok {
		return stmt.updateWithMap(dataMap, options...)
	}

	return stmt.updateWithStruct(data, options...)
}

func (stmt *UpdateStmt) updateWithStruct(data interface{}, options ...Option) error {
	var err error
	val := GetElem(data)
	valTye := val.Type()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("update data(%+v) is not a struct", data)
	}

	for i := 0; i < val.NumField(); i++ {
		dbColumn := valTye.Field(i).Tag.Get("db")
		if dbColumn == "" {
			continue
		}
		dbValue := val.Field(i).Interface()
		err = stmt.set(dbColumn+" = ?", dbValue, options...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (stmt *UpdateStmt) updateWithMap(data map[string]interface{}, options ...Option) error {
	var err error
	for dbColumn, dbValue := range data {
		err = stmt.set(dbColumn+" = ?", dbValue, options...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (stmt *UpdateStmt) clone() *UpdateStmt {
	copyArgs := make([]interface{}, len(stmt.args))
	copy(copyArgs, stmt.args)
	return &UpdateStmt{
		Builder: strings.Builder{},
		table:   stmt.table,
		args:    copyArgs,
	}
}
