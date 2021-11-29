# Sqlxx

Sqlxx is extended package for sqlx, it provide a sql builder and log function.

## Install

```
go get -u github.com/vx416/sqlxx
```
## Example

### Get DB Instance

```go
master, err := sqlx.Connect(
	"mysql",
	fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"test",
		"test",
		"localhost",
		"3306",
		"test_db",
	),
)
dao := sqlxx.NewWith(master)

slave, err := sqlx.Connect(
	"mysql",
	fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"test",
		"test",
		"localhost",
		"3306",
		"test_db",
	),
)

dao := sqlxx.NewWithCluster([]*sqlx.DB{master}, []*sqlx.DB{slave})
```

### SQL Builder

```go
ctx := context.Background()

dao := sqlxx.NewWithCluster([]*sqlx.DB{master}, []*sqlx.DB{slave})
users := []model.User{}
q := builder.Query().Select("id").From("users").
		And("id = ?", 1).AndIn("status IN (?)", []int{1,2,3})
db := dao.GetDB(ctx)

err := db.Select(ctx, &users, q)

type ListUsersOpt struct {
	ID           uint64        `sql:"col:id"`
	CreatedAtGte time.Time     `sql:"col:created_at;op:>="`
	NameLike     string        `sql:"col:name; op:%{}%"`
	StatusNotIn  []uint8       `sql:"col:status; op:notin"`
}

opt := ListUsersOpt{ID: 1, CreatedAtGte: time.Now(), NameLike: "vic"}
q := builder.Query().Select("id").From("users").Where(opt, builder.SkipZero)

```

[more examples](./builder/query_test.go)
