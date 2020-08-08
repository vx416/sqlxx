package sqlxx

import (
	"github.com/jmoiron/sqlx"
	"github.com/vicxu416/sqlxx/sqlbuilder"
)

type Insert struct {
}

func (queryCtx Insert) SetupContext(exec *Executor) error {
	insertStmt, err := sqlbuilder.Insert(exec.table, exec.data)
	if err != nil {
		return err
	}
	query := insertStmt.Sql(exec.db.DBType())
	exec.addQueryArgs(query, exec.data)
	return nil
}

func (queryCtx Insert) QueryHandle(exec *Executor, query string, args ...interface{}) error {
	switch exec.db.dbType {
	case PG:
		return queryCtx.pgQuery(exec, query, args[0])
	default:
		result, err := exec.db.NamedExec(query, args[0])
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		exec.LastInsertID = id
		return nil
	}
}

func (queryCtx Insert) pgQuery(exec *Executor, query string, arg interface{}) error {
	q, args, err := sqlx.BindNamed(sqlx.BindType(exec.db.DBType()), query, arg)
	if err != nil {
		return err
	}
	var id int64
	if err := exec.db.QueryRowx(q, args...).Scan(&id); err != nil {
		return err
	}
	exec.LastInsertID = id
	return nil
}

type BulkInsert struct {
}

func (queryCtx BulkInsert) SetupContext(exec *Executor) error {
	insertStmt, values, err := sqlbuilder.BulkInsert(exec.table, exec.data)
	if err != nil {
		return err
	}
	query := insertStmt.Sql(exec.db.DBType())
	exec.addQueryArgs(query, values...)
	return nil
}

func (queryCtx BulkInsert) QueryHandle(exec *Executor, query string, args ...interface{}) error {
	result, err := exec.db.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	exec.RowsAffected = rowsAffected

	return nil
}
