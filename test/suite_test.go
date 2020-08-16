package test

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vicxu416/sqlxx"
)

func TestCRUDSuite(t *testing.T) {
	suite.Run(t, new(CRUDSuite))
}

func InitDB(db *sqlxx.DB, file string) error {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	initSQL, err := ioutil.ReadFile(filepath.Join(dir, "./testdata/"+file+".sql"))
	if err != nil {
		return err
	}
	_, err = db.Exec(string(initSQL))

	return err
}

type CRUDSuite struct {
	suite.Suite
	db *sqlxx.DB
}

func (suite *CRUDSuite) SetupSuite() {
	db, err := sqlxx.Open("postgres", "user=test password=test host=127.0.0.1 port=5432 dbname=test sslmode=disable")
	suite.Require().Nil(err)
	suite.db = db
	err = InitDB(suite.db, "pg_init")
	suite.Require().Nil(err)
}
