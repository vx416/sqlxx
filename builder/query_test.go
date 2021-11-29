package builder

import (
	"database/sql"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vicxu416/sqlxx/logger"
	"gopkg.in/guregu/null.v4"
)

type TestCase struct {
	ExpectSql string
	argsLen   int
	builder   Builder
}

func (tc *TestCase) T(t *testing.T) {
	sql, args, err := tc.builder.Build()
	require.NoError(t, err)
	assert.Len(t, args, tc.argsLen)
	assert.Equal(t, tc.ExpectSql, logger.ExplainSQL(sql, args...))
}

func TestQuery_Select(t *testing.T) {
	tcs := []TestCase{
		{
			"SELECT * FROM users", 0, Query().From("users"),
		},
		{
			"SELECT id FROM users", 0, Query().From("users").Select("id"),
		},
		{
			"SELECT id, name FROM users", 0, Query().From("users").Select("id, name"),
		},
		{
			"SELECT id, name, gender as sex FROM users", 0,
			Query().From("users").Select("id, name").Select("gender as sex"),
		},
		{
			"SELECT id, name, gender as sex, phone FROM users", 0,
			Query().From("users").Select("id, name").Select("gender as sex, phone"),
		},
		{
			"SELECT id, name, gender as sex, phone, address, email FROM users", 0,
			Query().From("users").Select("id, name").Select("gender as sex, phone").Select([]string{"address", "email"}...),
		},
	}

	for _, tc := range tcs {
		t.Run("select", func(t *testing.T) {
			tc.T(t)
		})
	}
}

func TestQuery_SelectWhere(t *testing.T) {
	tcs := []TestCase{
		{
			"SELECT id, name, phone FROM users WHERE id = 1 OR gender IN (1, 2)", 3,
			Query().From("users").And("id = ?", 1).OrIn("gender IN (?)", []uint8{1, 2}).Select([]string{"id", "name", "phone"}...),
		},
		{
			"SELECT * FROM users WHERE id = 1 AND gender IN (3, 4, 1, 2)", 5,
			Query().From("users").And("id = ?", 1).AndIn("gender IN (?)", [][]int{{3, 4}, {1, 2}}),
		},
		{
			"SELECT * FROM users WHERE kind = 0 AND name LIKE \"%vic%\"", 1,
			Query().From("users").And("kind = 0", nil).And("name LIKE ?", "%vic%").And("id = ?", 0, SkipZero),
		},
		{
			"SELECT id, kind, name FROM users", 0,
			Query().From("users").And("id = ?", 0, SkipZero).And("kind IN (?)", []int{}, SkipZero).
				And("name = ?", "", SkipZero).And("decimal = ?", decimal.Zero, SkipZero).Select("id, kind, name"),
		},
	}

	for _, tc := range tcs {
		t.Run("selectWhere", func(t *testing.T) {
			tc.T(t)
		})
	}
}

func TestQuery_SelectJoinWhere(t *testing.T) {
	tcs := []TestCase{
		{
			"SELECT id, name, phone FROM users u JOIN projects p ON u.id = p.user_id WHERE id = 1 OR gender IN (1, 2)", 3,
			Query().From("users u").And("id = ?", 1).OrIn("gender IN (?)", []uint8{1, 2}).
				Select([]string{"id", "name", "phone"}...).Join("projects p ON u.id = p.user_id"),
		},
	}

	for _, tc := range tcs {
		t.Run("selectJoinWhere", func(t *testing.T) {
			tc.T(t)
		})
	}
}

type UserQuery struct {
	ID           uint64        `sql:"col:id"`
	Email        string        `sql:"col:email"`
	KindIn       []uint8       `sql:"col:kind;op:in"`
	CreatedAtGte time.Time     `sql:"col:created_at;op:>="`
	NameLike     string        `sql:"col:name; op:%{}%"`
	EmailLike    string        `sql:"col:email; op:{}%"`
	Status       sql.NullInt32 `sql:"col:status"`
	Kind         null.Int      `sql:"col:kind"`
	StatusNotIn  []uint8       `sql:"col:status; op:notin"`
}

func (user UserQuery) TableName() string {
	return "users"
}

func TestQuery_WhereStruct(t *testing.T) {
	tcs := []TestCase{
		{
			"SELECT * FROM users WHERE id = 1", 1,
			Query().Where(UserQuery{ID: 1}, SkipZero),
		},
		{
			"SELECT * FROM users WHERE id = 1 AND email = \"vic@gmail.com\"", 2,
			Query().Where(UserQuery{ID: 1, Email: "vic@gmail.com"}, SkipZero),
		},
		{
			"SELECT * FROM users u WHERE id = 1 AND kind IN (1, 2)", 3,
			Query().From("users u").Where(UserQuery{ID: 1, KindIn: []uint8{1, 2}}, SkipZero),
		},
		{
			"SELECT * FROM users u WHERE id = 1 AND created_at >= \"2021-07-31 12:30:03\"", 2,
			Query().From("users u").Where(UserQuery{ID: 1, CreatedAtGte: time.Date(2021, 7, 31, 12, 30, 3, 0, time.UTC)}, SkipZero),
		},
		{
			"SELECT * FROM users WHERE id = 1 AND name LIKE \"%vic%\"", 2,
			Query().Where(UserQuery{ID: 1, NameLike: "vic"}, SkipZero),
		},
		{
			"SELECT * FROM profiles WHERE id = 1 AND email LIKE \"vic%\"", 2,
			Query().From("profiles").Where(UserQuery{ID: 1, EmailLike: "vic"}, SkipZero),
		},
		{
			"SELECT * FROM profiles WHERE status = 0", 1,
			Query().From("profiles").Where(UserQuery{Status: sql.NullInt32{Valid: true}}, SkipZero),
		},
		{
			"SELECT * FROM profiles WHERE status NOT IN (1, 2)", 2,
			Query().From("profiles").Where(UserQuery{StatusNotIn: []uint8{1, 2}}, SkipZero),
		},
	}

	for _, tc := range tcs {
		t.Run("whereStruct", func(t *testing.T) {
			tc.T(t)
		})
	}
}

func TestQuery_SubQuery(t *testing.T) {
	tcs := []TestCase{
		{
			"SELECT * FROM ( SELECT AVG(amount) sum, user_id FROM account GROUP BY user_id ) stat WHERE user_id = 100", 1,
			Query().From("( ? ) stat", Query().From("account").Select("AVG(amount) sum, user_id").GroupBy("user_id")).And("user_id = ?", 100),
		},
		{
			"SELECT * FROM users WHERE id IN (SELECT user_id FROM profiles WHERE kind = 0) AND status IN (1, 2)", 3,
			Query().From("users").And("id IN (?)", Query().Select("user_id").From("profiles").And("kind = ?", 0)).AndIn("status IN (?)", []uint8{1, 2}),
		},
		{
			"SELECT * FROM ( SELECT AVG(amount) sum, user_id FROM account GROUP BY user_id ) stat WHERE user_id IN (SELECT user_id FROM profiles WHERE kind = 0) AND status IN (1, 2)", 3,
			Query().From("( ? ) stat", Query().From("account").Select("AVG(amount) sum, user_id").GroupBy("user_id")).
				And("user_id IN (?)", Query().Select("user_id").From("profiles").And("kind = ?", 0)).AndIn("status IN (?)", []uint8{1, 2}),
		},
		{
			"SELECT * FROM users JOIN (SELECT id FROM users LIMIT 100 OFFSET 10000) u1 ON users.id = u1.id WHERE status = 1", 1,
			Query().From("users").Join("(?) u1 ON users.id = u1.id", Query().From("users").Select("id").LimitOffset(100, 10000)).
				And("status = ?", 1),
		},
	}

	for _, tc := range tcs {
		t.Run("subQuery", func(t *testing.T) {
			tc.T(t)
		})
	}
}

func TestQuery_Other(t *testing.T) {
	tcs := []TestCase{
		{
			"SELECT * FROM users ORDER BY id DESC, created_at ASC", 0,
			Query().From("users").OrderBy("id DESC", "created_at ASC"),
		},
		{
			"SELECT id, status FROM users GROUP BY id, status ORDER BY id DESC, created_at ASC", 0,
			Query().Select("id, status").From("users").GroupBy("id").OrderBy("id DESC", "created_at ASC").GroupBy("status"),
		},
		{
			"SELECT id, status FROM users GROUP BY id, status LIMIT 100 OFFSET 1", 0,
			Query().Select("id, status").From("users").GroupBy("id").LimitOffset(100, 1).GroupBy("status"),
		},
		{
			"SELECT id, status FROM users GROUP BY id, status LIMIT 100 OFFSET 1 FOR UPDATE", 0,
			Query().Select("id, status").From("users").GroupBy("id").LimitOffset(100, 1).GroupBy("status").Lock(string(MSWRITELOCK)),
		},
	}

	for _, tc := range tcs {
		t.Run("other", func(t *testing.T) {
			tc.T(t)
		})
	}
}

func TestQuery_Union(t *testing.T) {
	tcs := []TestCase{
		{
			"(SELECT * FROM users WHERE s_user_id = 1) UNION ALL (SELECT * FROM users WHERE b_user_id = 1)", 2,
			Query().From("users").And("s_user_id = ?", 1).UnionAll(Query().From("users").And("b_user_id = ?", 1)),
		},
		{
			"(SELECT * FROM users WHERE s_user_id = 1) UNION (SELECT * FROM users WHERE b_user_id = 1) UNION (SELECT * FROM users WHERE b_user_id = 2)", 3,
			Query().From("users").And("s_user_id = ?", 1).Union(Query().From("users").And("b_user_id = ?", 1)).Union(Query().From("users").And("b_user_id = ?", 2)),
		},
		{
			"(SELECT * FROM users WHERE s_user_id = 1 ORDER BY id LIMIT 100) UNION ALL (SELECT * FROM users WHERE b_user_id = 1 ORDER BY id LIMIT 100)", 2,
			Query().From("users").And("s_user_id = ?", 1).OrderBy("id").LimitOffset(100, 0).UnionAll(Query().From("users").And("b_user_id = ?", 1).OrderBy("id").LimitOffset(100, 0)),
		},
	}

	for _, tc := range tcs {
		t.Run("other", func(t *testing.T) {
			tc.T(t)
		})
	}
}
