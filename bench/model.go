package bench

import "time"

type User struct {
	ID           int64  `gorm:"column:id" db:"id"`
	FirstName    string `gorm:"column:first_name" db:"first_name"`
	LastName     string `gorm:"column:last_name" db:"last_name"`
	PasswordHash []byte `gorm:"column:password_hash" db:"password_hash"`
	CreatedAt    string `gorm:"column:created_at" db:"created_at"`
}

func NewUser() User {
	return User{
		FirstName:    "test1234567",
		LastName:     "test1234567",
		PasswordHash: []byte("7y12hijnwfliasyhoi243poiasfjh2u3i4yhasujf"),
		CreatedAt:    time.Now().Format(time.RFC1123),
	}
}
