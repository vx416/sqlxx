package test

import (
	"context"
	"time"

	"github.com/vicxu416/sqlxx"
)

func (suite *CRUDSuite) TestTransactionLock() {
	user := NewUser()
	exec := suite.db.Insert("users", user)
	err := exec.Do()
	suite.Require().NoError(err)
	id := exec.LastInsertID()

	query := sqlxx.NewQueryOpts()
	query.Lock("UPDATE")
	query.Where("id = ?", id)

	tx, err := suite.db.Begin()
	suite.Require().NoError(err)

	users := []*User{}
	err = tx.Select("users", &users, query).Do()
	suite.Require().NoError(err)

	done := make(chan struct{})

	go func() {
		err = suite.db.Select("users", &users, query).Do()
		suite.Require().NoError(err)
		close(done)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	i := 0

	select {
	case <-done:
		i = 1
	case <-ctx.Done():
		i = 2
	}

	suite.Assert().Equal(i, 2)
}

func (suite *CRUDSuite) TestTransactionRollback() {
	user := NewUser()
	tx, err := suite.db.Begin()
	suite.Require().NoError(err)
	exec := tx.Insert("users", user)
	err = exec.Do()
	suite.Require().NoError(err)
	id := exec.LastInsertID()
	err = tx.Rollback()
	suite.Require().NoError(err)

	query := sqlxx.NewQueryOpts()
	query.Where("id = ?", id)
	users := []*User{}
	err = suite.db.Select("users", &users, query).Do()
	suite.Require().NoError(err)
	suite.Assert().Len(users, 0)
}
