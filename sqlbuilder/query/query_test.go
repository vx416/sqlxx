package query

type temp struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// func TestWhere(t *testing.T) {
// 	data := temp{
// 		ID:   1,
// 		Name: "test",
// 	}

// 	where := NewWhere()
// 	where.Where("id = ?", 1)
// 	where.Where("first_name = ?", "test")
// 	where.AndStruct(data)
// 	where.In("id", []int32{1, 2, 3, 4, 4})
// 	where.Or("last_name = ?", "hello")
// 	t.Logf("QUERY:%s", where.Query())
// 	t.Logf("ARGS:%+v", where.Args())
// 	if where.lastErr != nil {
// 		t.Logf("ERR:%+v", where.lastErr)
// 	}
// }

// func TestSelect(t *testing.T) {
// 	sel := NewSelect()

// 	sel.Select("id, name").Select("count(*) as t").
// 		Group("id, name").Having("t > 2").Group("t2").
// 		Limit(10).Offset(10).Orderby("id", "DESC")

// 	queries := sel.Queries()
// 	t.Logf("options: %s", queries["options"])
// }
