package test

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64        `db:"id"`
	FirstName    string       `db:"first_name"`
	LastName     string       `db:"last_name"`
	PasswordHash []byte       `db:"password_hash"`
	CreatedAt    sql.NullTime `db:"created_at"`
}

func NewUser() User {
	return User{
		FirstName:    "test1234567",
		LastName:     "test1234567",
		PasswordHash: []byte("1239asf82394890zdgklmc 09234ijoasdi90ser"),
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}
}
