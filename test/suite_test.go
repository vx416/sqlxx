package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vicxu416/sqlxx"
	"github.com/vicxu416/sqlxx/testdata"
)

func TestCRUDSuite(t *testing.T) {
	suite.Run(t, new(CRUDSuite))
}

type CRUDSuite struct {
	suite.Suite
	db *sqlxx.DB
}

func (suite *CRUDSuite) SetupSuite() {
	db, err := sqlxx.Open("postgres", "user=dev password=dev host=127.0.0.1 port=5432 dbname=test sslmode=disable")
	// db, err := sqlxx.Open("sqlite3", ":memory:")
	suite.Require().Nil(err)
	suite.db = db
	err = testdata.InitDB(suite.db, "pg_init")
	suite.Require().Nil(err)
}
