package bench

import (
	"testing"

	"github.com/vicxu416/sqlxx/testdata"
)

func BenchmarkInsert(b *testing.B) {
	b.Run("sqlxx", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			user := testdata.NewUser()

			if err := sqlxxDB.Insert("users", user).Do(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := testdata.NewUser()

			err := gormDB.Create(&user).Error
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := testdata.NewUser()

			err := gorpDB.Insert(&user)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
