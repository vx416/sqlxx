package test

import (
	"github.com/vicxu416/sqlxx"
)

func (suite *CRUDSuite) TestSelectAll() {
	user := NewUser()
	err := suite.db.Insert("users", &user).Do()
	suite.Require().Nil(err)
	users := []*User{}
	exec := suite.db.Select("users", &users, nil)
	err = exec.Do()
	suite.Require().Nil(err)
	suite.Require().Greater(len(users), 0)
	suite.Assert().NotZero(users[0].ID)
	suite.Assert().NotEmpty(users[0].FirstName)
	suite.Assert().NotEmpty(users[0].LastName)
}

func (suite *CRUDSuite) TestSelectWhere() {
	user := NewUser()
	exec := suite.db.Insert("users", &user)
	err := exec.Do()
	suite.Require().Nil(err)

	query := sqlxx.NewQueryOpts()
	query.Select("id").In("id", []int64{1, 2, 3, 4, 5, 6, exec.LastInsertID()}).Where("first_name <> ?", "")

	users := []*User{}
	exec = suite.db.Select("users", &users, query)
	err = exec.Do()
	suite.Require().Nil(err)
	suite.Require().Greater(len(users), 0)
	suite.Assert().NotZero(users[0].ID)
	suite.Assert().Empty(users[0].FirstName)
}

func (suite *CRUDSuite) TestGet() {
	user := NewUser()
	exec := suite.db.Insert("users", &user)
	err := exec.Do()
	suite.Require().Nil(err)

	query := sqlxx.NewQueryOpts()
	query.Where("id = ?", user.ID)
	user2 := User{}
	err = suite.db.Get("users", &user2, query).Do()
	suite.Require().Nil(err)
	suite.Equal(user2.ID, user.ID)
}
