package testdata

import (
	"database/sql"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"runtime"
	"strconv"

	"time"
)

type exec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func InitDB(db exec, file string) error {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	initSQL, err := ioutil.ReadFile(filepath.Join(dir, "./"+file+".sql"))
	if err != nil {
		return err
	}
	_, err = db.Exec(string(initSQL))

	return err
}

func GenPWD() []byte {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(1000)
	return []byte(strconv.Itoa(i))
}

type User struct {
	ID           int64  `gorm:"column:id" db:"id"`
	FirstName    string `gorm:"column:first_name" db:"first_name"`
	LastName     string `gorm:"column:last_name" db:"last_name"`
	PasswordHash []byte `gorm:"column:password_hash" db:"password_hash"`
	// CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func NewUser() User {
	return User{
		FirstName:    "test1234567",
		LastName:     "test1234567",
		PasswordHash: GenPWD(),
		// CreatedAt:    time.Now(),
	}
}
