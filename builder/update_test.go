package builder

import (
	"testing"

	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
)

func TestUpdate_Where(t *testing.T) {
	tcs := []TestCase{
		{
			"UPDATE users SET name = \"vic\" WHERE id = 1", 2,
			Update().Table("users").Set("name = ?", "vic").And("id = ?", 1),
		},
		{
			"UPDATE users SET phone = 1234, name = \"haha\" WHERE status IN (1, 2, 3)", 4,
			Update().AndIn("status IN (?)", []int{1, 2, 3}).Table("users").Set("phone = 1234", nil).Set("name = ?", "haha"),
		},
		// {
		// 	"UPDATE users SET level = 1 WHERE id = 1", 2,
		// 	Update().Where(UserQuery{ID: 1}, SkipZero).Set("level = ?", 1).Set("status = ?", 0, SkipZero),
		// },
		{
			"UPDATE users SET level = 1", 1,
			Update().Table("users").Set("level = ?", 1).Set("status = ?", 0, SkipZero),
		},
	}

	for _, tc := range tcs {
		t.Run("update", func(t *testing.T) {
			tc.T(t)
		})
	}
}

type User struct {
	ID        int                 `db:"id"`
	Name      string              `db:"name"`
	Level     int                 `db:"level"`
	Status    uint8               `db:"status"`
	CreatedAt null.Time           `db:"created_at"`
	Money     decimal.NullDecimal `db:"money"`
}

func (t User) TableName() string {
	return "users"
}

func TestUpdate_With(t *testing.T) {
	tcs := []TestCase{
		{
			"UPDATE users SET name = \"vic\", level = 2 WHERE id = 1", 3,
			Update().SetWith(&User{Name: "vic", Level: 2}, SkipZero).And("id = ?", 1),
		},
		{
			`UPDATE users SET status = 2, money = "100" WHERE id = 1`, 3,
			Update().SetWith(&User{Status: 2, Money: decimal.NullDecimal{Valid: true, Decimal: decimal.NewFromInt(100)}}, SkipZero).And("id = ?", 1),
		},
		{
			`UPDATE users SET level = 1, id = 2 WHERE id = 1`, 3,
			Update().Table("users").SetWith(map[string]interface{}{"level": 1, "id": 2}).And("id = ?", 1),
		},
	}

	for _, tc := range tcs {
		t.Run("updateWith", func(t *testing.T) {
			tc.T(t)
		})
	}
}
