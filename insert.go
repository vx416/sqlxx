package sqlxx

import (
	"reflect"
)

type Insert struct {
}

func (execCtx Insert) SetupContext(exec *Executor) (query string, args []interface{}, err error) {
	query, data, err := exec.db.GetBuilder().BuildInsert(exec.table, exec.data)
	if err != nil {
		return "", nil, err
	}
	return query, data, nil
}

func (execCtx Insert) ExecHandle(exec *Executor, query string, args ...interface{}) error {
	switch exec.db.driverName {
	case PG:
		err := execCtx.pgQuery(exec, query, args...)
		if err != nil {
			return err
		}
	default:
		result, err := exec.db.Exec(query, args...)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		exec.lastInsertID = id
	}

	val := reflect.ValueOf(exec.data)
	if val.Kind() == reflect.Ptr {
		valElem := val.Elem()
		for i := 0; i < valElem.NumField(); i++ {
			if valElem.Type().Field(i).Tag.Get("db") == "id" {
				valElem.Field(i).SetInt(exec.lastInsertID)
			}
		}
	}

	return nil
}

func (execCtx Insert) pgQuery(exec *Executor, query string, args ...interface{}) error {
	var id int64
	if err := exec.db.QueryRowx(query, args...).Scan(&id); err != nil {
		return err
	}
	exec.lastInsertID = id
	return nil
}

type BulkInsert struct {
}

func (execCtx BulkInsert) SetupContext(exec *Executor) (string, []interface{}, error) {
	query, values, err := exec.db.GetBuilder().BuildBulkInsert(exec.table, exec.data)
	if err != nil {
		return "", nil, err
	}
	return query, values, nil
}

func (execCtx BulkInsert) ExecHandle(exec *Executor, query string, args ...interface{}) error {
	result, err := exec.db.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	exec.rowsAffected = rowsAffected

	return nil
}
