package builder

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestInsert(t *testing.T) {
	tcs := []TestCase{
		{
			"INSERT INTO users (id, user_id, name) VALUES (1, 1, \"vic\")", 3,
			Insert().Table("users").InsertRows(map[string]interface{}{"id": 1, "user_id": 1, "name": "vic"}),
		},
		{
			"INSERT INTO users (name, level, status, money) VALUES (\"vic\", 1, 1, \"100\")", 4,
			Insert().InsertRows(User{Name: "vic", Level: 1, Status: 1, Money: decimal.NewNullDecimal(decimal.New(100, 0))}),
		},
		{
			"INSERT INTO users (name, level, status, money) VALUES (\"vic\", 1, 1, \"100\"), (\"vic2\", 2, 2, \"200\")", 8,
			Insert().InsertRows([]User{{Name: "vic", Level: 1, Status: 1, Money: decimal.NewNullDecimal(decimal.New(100, 0))},
				{Name: "vic2", Level: 2, Status: 2, Money: decimal.NewNullDecimal(decimal.New(200, 0))}}),
		},
		{
			"INSERT INTO users (id, user_id, name) VALUES (1, 1, \"vic\"), (2, 2, \"vic2\")", 6,
			Insert().Table("users").InsertRows([]map[string]interface{}{{"id": 1, "user_id": 1, "name": "vic"}, {"id": 2, "user_id": 2, "name": "vic2"}}),
		},
	}

	for _, tc := range tcs {
		t.Run("insert", func(t *testing.T) {
			tc.T(t)
		})
	}
}
