package test

func (suite *CRUDSuite) TestInsert() {
	user := NewUser()
	exec := suite.db.Insert("users", &user)
	err := exec.Do()
	suite.Require().Nil(err)
	suite.Assert().NotZero(user.ID)
	suite.Assert().Equal(user.ID, exec.LastInsertID())
}

func (suite *CRUDSuite) TestBulkInsert() {
	users := make([]*User, 0, 10)
	for i := 0; i < 10; i++ {
		user := NewUser()
		users = append(users, &user)
	}
	exec := suite.db.BulkInsert("users", users)
	err := exec.Do()
	suite.Require().Nil(err)
	suite.Assert().NotZero(exec.RowsAffected())
	suite.Assert().Equal(exec.RowsAffected(), int64(10))
}
