package sqlxx

import (
	"reflect"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

type DBType string

const (
	PG DBType = "postgres"
)

type DB struct {
	*proxyDB

	dbType DBType
}

func Open(driver, dataSourceName string) (*DB, error) {
	db, err := sqlx.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	proxyDB := newProxyDB(db)
	return &DB{proxyDB: proxyDB, dbType: DBType(driver)}, nil
}

func (db *DB) Insert(table string, data interface{}) *Executor {
	var queryContext QueryContext = Insert{}
	if reflect.TypeOf(data).Kind() == reflect.Slice {
		queryContext = BulkInsert{}
	}

	return NewExec(queryContext, db, table, data)
}

func (db *DB) Debug(on bool) {
	db.proxyDB.debug = on
}

func (db DB) DBType() string {
	return string(db.dbType)
}
