package sqlxx

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

func newProxyDB(dbx *sqlx.DB) *proxyDB {
	return &proxyDB{
		DB:     dbx,
		logger: newDefualtLogger(),
		debug:  true,
	}

}

type proxyDB struct {
	DB     *sqlx.DB
	Tx     *sqlx.Tx
	logger Logger
	debug  bool
}

func (db *proxyDB) Clone() *proxyDB {
	return &proxyDB{
		DB:     db.DB,
		Tx:     db.Tx,
		logger: db.logger,
		debug:  db.debug,
	}
}

func (db *proxyDB) begin() (*proxyDB, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	txProxy := db.Clone()
	txProxy.Tx = tx

	return txProxy, nil
}

func (db *proxyDB) rollback() error {
	return db.Tx.Rollback()
}

func (db *proxyDB) commit() error {
	return db.Tx.Commit()
}

func (db *proxyDB) Logf(msg string, args ...interface{}) {
	if db.debug {
		db.logger.Logf(msg, args...)
	}
}

func (db *proxyDB) DriverName() string {
	return db.DB.DriverName()
}

func (db *proxyDB) Rebind(query string) string {
	return db.DB.Rebind(query)
}

func (db *proxyDB) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return db.DB.BindNamed(query, arg)
}

func (db *proxyDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	db.Logf(query, args...)

	if db.Tx != nil {
		return db.Tx.Exec(query, args...)
	}

	return db.DB.Exec(query, args...)
}

func (db *proxyDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	db.Logf(query, args...)

	if db.Tx != nil {
		return db.Tx.Query(query, args...)
	}

	return db.DB.Query(query, args...)
}

func (db *proxyDB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	db.Logf(query, args...)

	if db.Tx != nil {
		return db.Tx.Queryx(query, args...)
	}

	return db.DB.Queryx(query, args...)
}

func (db *proxyDB) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	db.Logf(query, args...)

	if db.Tx != nil {
		return db.Tx.QueryRowx(query, args...)
	}

	return db.DB.QueryRowx(query, args...)
}

func (db *proxyDB) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return sqlx.NamedQuery(db, query, arg)
}

func (db *proxyDB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return sqlx.NamedExec(db, query, arg)
}

func (db *proxyDB) Select(dest interface{}, query string, args ...interface{}) error {
	return sqlx.Select(db, dest, query, args...)
}

func (db *proxyDB) Get(dest interface{}, query string, args ...interface{}) error {
	return sqlx.Get(db, dest, query, args...)
}

func newDefualtLogger() Logger {
	return logger{Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime)}
}

type Logger interface {
	Logf(msg string, args ...interface{})
}

type logger struct {
	*log.Logger
}

func (l logger) Logf(msg string, args ...interface{}) {
	var logMsg strings.Builder
	logMsg.WriteString(colorize("\nSQL---\n", colorYellow))
	logMsg.WriteString(colorize(msg, colorGreen))
	if len(args) > 0 {
		logMsg.WriteRune('\n')
		logMsg.WriteString(colorize(args, colorBlue))
	}
	logMsg.WriteString(colorize("\nSQL---", colorYellow))
	l.Logger.Printf(logMsg.String())
}

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
