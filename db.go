package sqlxx

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vicxu416/sqlxx/builder"
	"github.com/vicxu416/sqlxx/logger"
)

var ErrNilTx = errors.New("tx is nil")

type DB struct {
	Cluster *Cluster
	Tx      *sqlx.Tx
}

func (db *DB) GetRawDB(ctx context.Context) (*sql.DB, error) {
	sqlxDB, err := db.Cluster.GetDB(ctx)
	if err != nil {
		return nil, err
	}
	return sqlxDB.DB, nil
}

func (db *DB) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	exec, err := db.getExec(ctx)
	if err != nil {
		return nil, err
	}
	return sqlx.NamedQueryContext(ctx, exec, query, arg)
}

func (db *DB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	exec, err := db.getExec(ctx)
	if err != nil {
		return nil, err
	}
	return sqlx.NamedExecContext(ctx, exec, query, arg)
}

func (db *DB) Select(ctx context.Context, dest interface{}, query builder.Builder) error {
	queryS, args, err := query.Build()
	if err != nil {
		return err
	}
	return db.SelectContext(ctx, dest, queryS, args...)
}

func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var (
		start  = time.Now()
		err    error
		driver sqlx.QueryerContext
	)
	defer func() {
		cost := time.Since(start)
		logger.Print(ctx, 0, err, cost, query, args...)
	}()

	driver, err = db.getQuery(ctx)
	if err != nil {
		return err
	}
	err = sqlx.SelectContext(ctx, driver, dest, query, args...)
	return err
}

func (db *DB) Get(ctx context.Context, dest interface{}, query builder.Builder) error {
	queryS, args, err := query.Build()
	if err != nil {
		return err
	}
	return db.GetContext(ctx, dest, queryS, args...)
}

func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var (
		start  = time.Now()
		err    error
		driver sqlx.QueryerContext
	)
	defer func() {
		cost := time.Since(start)
		logger.Print(ctx, 0, err, cost, query, args...)
	}()

	driver, err = db.getQuery(ctx)
	if err != nil {
		return err
	}
	err = sqlx.GetContext(ctx, driver, dest, query, args...)
	return err
}

func (db *DB) Exec(ctx context.Context, query builder.Builder) (sql.Result, error) {
	queryS, args, err := query.Build()
	if err != nil {
		return nil, err
	}
	return db.ExecContext(ctx, queryS, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var (
		start = time.Now()
		err   error
		res   sql.Result
		rows  int64
	)
	defer func() {
		cost := time.Since(start)
		logger.Print(ctx, rows, err, cost, query, args...)
	}()

	if db.Tx != nil {
		res, err = db.Tx.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		rows, err = res.RowsAffected()
		if err != nil {
			return res, err
		}
		return res, err
	}

	sqlxDB, err := db.Cluster.GetDB(ctx)
	if err != nil {
		return nil, err
	}
	res, err = sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return res, err
	}
	rows, err = res.RowsAffected()
	if err != nil {
		return res, err
	}
	return res, err
}

func (db *DB) Begin(ctx context.Context, txOpt *sql.TxOptions) (*DB, error) {
	if txOpt == nil {
		txOpt = &sql.TxOptions{}
	}
	if IsReadOnly(ctx) {
		txOpt.ReadOnly = true
	}

	sqlxDB, err := db.Cluster.GetDB(ctx)
	if err != nil {
		return nil, err
	}
	sqlxTx, err := sqlxDB.BeginTxx(ctx, txOpt)
	if err != nil {
		return nil, err
	}
	return &DB{Tx: sqlxTx, Cluster: nil}, nil
}

func (db *DB) Commit(ctx context.Context) error {
	if db.Tx == nil {
		return ErrNilTx
	}

	return db.Tx.Commit()
}

func (db *DB) Rollback(ctx context.Context) error {
	if db.Tx == nil {
		return ErrNilTx
	}

	return db.Tx.Rollback()
}

func (db *DB) IsTx() bool {
	return db.Tx != nil
}

func (db *DB) getExec(ctx context.Context) (sqlx.ExtContext, error) {
	if db.Tx != nil {
		return SqlxxExtContext{db.Tx}, nil
	}

	sqlxDB, err := db.Cluster.GetDB(ctx)
	if err != nil {
		return nil, err
	}
	return SqlxxExtContext{sqlxDB}, nil
}

func (db *DB) getQuery(ctx context.Context) (sqlx.QueryerContext, error) {
	if db.Tx != nil {
		return db.Tx, nil
	}

	if !IsMaster(ctx) {
		ctx = WithSlave(ctx)
	}
	sqlxDB, err := db.Cluster.GetDB(ctx)
	if err != nil {
		return nil, err
	}
	return sqlxDB, nil
}

type SqlxxExtContext struct {
	sqlx.ExtContext
}

func (exec SqlxxExtContext) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var (
		start = time.Now()
		err   error
		res   sql.Result
		rows  int64
	)
	defer func() {
		cost := time.Since(start)
		logger.Print(ctx, rows, err, cost, query, args...)
	}()

	res, err = exec.ExtContext.ExecContext(ctx, query, args...)
	if err != nil {
		return res, err
	}
	rows, err = res.RowsAffected()
	if err != nil {
		return res, err
	}
	return res, err
}

func (exec SqlxxExtContext) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var (
		start = time.Now()
		err   error
		res   *sql.Rows
	)
	defer func() {
		cost := time.Since(start)
		logger.Print(ctx, 0, err, cost, query, args...)
	}()

	res, err = exec.ExtContext.QueryContext(ctx, query, args...)
	return res, err
}
