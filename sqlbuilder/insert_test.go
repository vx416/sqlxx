package sqlbuilder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	type source struct {
		ID        int64     `db:"id"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"time"`
		Ignore    string
	}
	data := []source{{
		ID:        1,
		Name:      "hello",
		CreatedAt: time.Now(),
	}}

	stmt, err := Insert("test", data)
	assert.Nil(t, err)
	t.Log(stmt.Sql("postgres"))
}
