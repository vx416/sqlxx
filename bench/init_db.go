package bench

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/vicxu416/sqlxx/testdata"

	"gopkg.in/gorp.v1"

	"github.com/jinzhu/gorm"
	"github.com/vicxu416/sqlxx"
)

var sqlxxDB *sqlxx.DB
var gormDB *gorm.DB
var gorpDB *gorp.DbMap

func init() {
	var err error
	sqlxxDB, err = OpenSqlxx()
	sqlxxDB.Debug(false)
	if err != nil {
		panic(err)
	}
	err = initDB(sqlxxDB)
	if err != nil {
		panic(err)
	}
	gormDB, err = OpenGorm()
	if err != nil {
		panic(err)
	}
	err = initDB(gormDB.DB())
	if err != nil {
		panic(err)
	}
	gorpDB, err = OpenGorp()
	if err != nil {
		panic(err)
	}
	err = initDB(gorpDB.Db)
	if err != nil {
		panic(err)
	}
	gorpDB.AddTableWithName(testdata.User{}, "users").SetKeys(true, "id")

}

type exec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func initDB(db exec) error {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	initSQL, err := ioutil.ReadFile(filepath.Join(dir, "../testdata/init.sql"))
	if err != nil {
		return err
	}
	_, err = db.Exec(string(initSQL))

	return err
}

func OpenSqlxx() (*sqlxx.DB, error) {
	return sqlxx.Open("sqlite3", ":memory:")
}

func OpenGorm() (*gorm.DB, error) {
	return gorm.Open("sqlite3", ":memory:")
}

func OpenGorp() (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	return dbmap, nil
}
