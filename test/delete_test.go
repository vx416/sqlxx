package test

import (
	"github.com/vicxu416/sqlxx"
)

func (suite *CRUDSuite) TestDelete() {
	err := suite.db.Delete("users", nil).Do()
	suite.Require().NoError(err)

	user := NewUser()
	exec := suite.db.Insert("users", user)
	err = exec.Do()
	suite.Require().NoError(err)
	query := sqlxx.NewQueryOpts()
	query.Where("id = ?", exec.LastInsertID())
	err = suite.db.Delete("users", query).Do()
	suite.Require().NoError(err)

}
