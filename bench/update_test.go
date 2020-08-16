package bench

import (
	"testing"

	"github.com/vicxu416/sqlxx"
)

func BenchmarkUpdate(b *testing.B) {
	user := User{FirstName: "testupdate", LastName: "testupdate", CreatedAt: "testcreated", ID: 1}

	b.Run("sqlxx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			query := sqlxx.NewQueryOpts()
			query.Where("id = ?", 1)

			if err := sqlxxDB.Update("users", &user, query).Do(); err != nil {
				b.Fatal(err)
			}

		}
	})

	b.Run("sqlxx updateStruct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := sqlxxDB.UpdateStruct("users", &user).Do(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := gormDB.Table("users").Where("id = ?", 1).Update(&user).Error; err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := gorpDB.Update(&user)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
