package bench

import (
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	b.Run("sqlxx", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			user := NewUser()

			if err := sqlxxDB.Insert("users", user).Do(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := NewUser()

			err := gormDB.Create(&user).Error
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := NewUser()

			err := gorpDB.Insert(&user)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkBulkInsert(b *testing.B) {
	b.Run("sqlxx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			users := []User{NewUser(), NewUser()}
			if err := sqlxxDB.BulkInsert("users", users).Do(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for i := 0; i < 2; i++ {
				user := NewUser()
				err := gormDB.Create(&user).Error
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user1 := NewUser()
			user2 := NewUser()
			err := gorpDB.Insert(&user1, &user2)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
