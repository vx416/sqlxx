package bench

import (
	"testing"

	"github.com/vicxu416/sqlxx"
)

func BenchmarkSelectAll(b *testing.B) {
	b.Run("sqlxx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			users := []User{}
			if err := sqlxxDB.Select("users", &users, nil).Do(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			users := []User{}
			if err := gormDB.Table("users").Scan(&users).Error; err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			users := []User{}
			_, err := gorpDB.Select(&users, "SELECT * FROM users")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSelectComplex(b *testing.B) {
	b.Run("sqlxx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			query := sqlxx.NewQueryOpts()
			query.Select("id, first_name, last_name").Limit(1).Offset(1).Group("id")
			query.In("id", []int{1, 2, 3, 4, 5}).Where("id <> ?", 5)
			users := []User{}
			if err := sqlxxDB.Select("users", &users, query).Do(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			users := []User{}
			if err := gormDB.Table("users").Select("id, first_name, last_name").
				Where("id <> ?", 5).Where("id IN (?)", []int{1, 2, 3, 4, 5}).
				Group("id").Offset(1).Limit(1).Scan(&users).Error; err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			users := []User{}
			_, err := gorpDB.Select(&users, `SELECT id, first_name, last_name FROM users 
			Where id <> ? AND id IN (?, ?, ?, ?, ?) GROUP BY id LIMIT ? OFFSET ?`, 5, 1, 2, 3, 4, 5, 1, 1)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
