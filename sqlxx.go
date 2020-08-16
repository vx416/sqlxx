package sqlxx

import (
	"sync"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicxu416/sqlxx/sqlbuilder"

	"github.com/jmoiron/sqlx"
)

type DBDriver string

const (
	PG DBDriver = "postgres"
)

type DB struct {
	*proxyDB
	exec        *Executor
	driverName  DBDriver
	builderPool sync.Pool
}

func Open(driver, dataSourceName string) (*DB, error) {
	db, err := sqlx.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	builderPool := sync.Pool{
		New: func() interface{} {
			return sqlbuilder.New(driver)
		},
	}

	proxyDB := newProxyDB(db)
	return &DB{proxyDB: proxyDB, driverName: DBDriver(driver), builderPool: builderPool}, nil
}

func (db *DB) GetBuilder() *sqlbuilder.Builder {
	return db.builderPool.Get().(*sqlbuilder.Builder)
}

func (db *DB) Insert(table string, data interface{}) Execr {
	return NewExec(db, Insert{}, table, data)
}

func (db *DB) BulkInsert(table string, data interface{}) Execr {
	return NewExec(db, BulkInsert{}, table, data)
}

func (db *DB) Update(table string, data interface{}, options *QueryOptions) Execr {
	exec := NewExec(db, UpdateCtx{}, table, data)
	if options == nil {
		options = NewQueryOpts()
	}
	exec.QueryOpts = options
	return exec
}

func (db *DB) UpdateStruct(table string, data interface{}) Execr {
	return NewExec(db, UpdateStructCtx{}, table, data)
}

func (db *DB) Delete(table string, options *QueryOptions, usingTable ...string) Execr {
	exec := NewExec(db, DeleteCtx{}, table, usingTable)
	if options == nil {
		options = NewQueryOpts()
	}
	exec.QueryOpts = options
	return exec
}

func (db *DB) Select(table string, data interface{}, options *QueryOptions) Execr {
	exec := NewExec(db, SelectCtx{}, table, data)
	if options == nil {
		options = NewQueryOpts()
	}
	exec.QueryOpts = options
	return exec
}

func (db *DB) Debug(on bool) {
	db.proxyDB.debug = on
}

func (db *DB) Begin() (*DB, error) {
	txDB := db.Clone()
	txProxy, err := txDB.proxyDB.begin()
	if err != nil {
		return nil, err
	}
	txDB.proxyDB = txProxy
	return txDB, nil
}

func (db *DB) Commit() error {
	return db.proxyDB.commit()
}

func (db *DB) Rollback() error {
	return db.proxyDB.rollback()
}

func (db *DB) Clone() *DB {
	return &DB{
		proxyDB:     db.proxyDB,
		builderPool: db.builderPool,
		exec:        db.exec,
		driverName:  db.driverName,
	}
}

func Transaction(db *DB, callback func(tx *DB) error) error {
	var callbackErr error
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if callbackErr = callback(tx); callbackErr != nil {
		err = tx.proxyDB.rollback()
	} else {
		err = tx.proxyDB.commit()
	}

	if err != nil {
		return err
	}
	return callbackErr
}

func (db DB) DriverName() string {
	return string(db.driverName)
}
