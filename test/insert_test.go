package test

import "github.com/vicxu416/sqlxx/testdata"

func (suite *CRUDSuite) TestInsert() {
	user := testdata.NewUser()
	exec := suite.db.Insert("users", &user)
	err := exec.Do()
	suite.Require().Nil(err)
	suite.Assert().NotZero(exec.LastInsertID)
}

func (suite *CRUDSuite) TestBulkInsert() {
	users := make([]*testdata.User, 0, 10)
	for i := 0; i < 10; i++ {
		user := testdata.NewUser()
		users = append(users, &user)
	}
	exec := suite.db.Insert("users", users)
	err := exec.Do()
	suite.Require().Nil(err)
	suite.Assert().NotZero(exec.RowsAffected)
	suite.Assert().Equal(exec.RowsAffected, int64(10))
}
