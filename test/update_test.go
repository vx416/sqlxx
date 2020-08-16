package test

import (
	"database/sql"
	"time"

	"github.com/vicxu416/sqlxx"
)

func (suite *CRUDSuite) TestUpdate() {
	user := NewUser()
	user.FirstName = "testFirstName"
	exec := suite.db.Insert("users", user)
	err := exec.Do()
	suite.Require().NoError(err)
	user2 := User{}
	user2.ID = exec.LastInsertID()
	queryOpt := sqlxx.NewQueryOpts()
	queryOpt.AndStruct(user2)
	user2.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	user2.FirstName = "updatedName"

	err = suite.db.Update("users", user2, queryOpt).Do()
	suite.Require().NoError(err)
}

func (suite *CRUDSuite) TestUpdateStruct() {
	user := NewUser()
	user.FirstName = "testFirstName"
	exec := suite.db.Insert("users", user)
	err := exec.Do()
	suite.Require().NoError(err)
	user2 := User{}
	user2.ID = exec.LastInsertID()
	user2.FirstName = "updatedName"
	err = suite.db.UpdateStruct("users", user2).Do()
	suite.Require().NoError(err)
}
