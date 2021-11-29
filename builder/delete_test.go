package builder

import "testing"

func TestDelete(t *testing.T) {
	tcs := []TestCase{
		{
			"DELETE FROM users WHERE id = 1", 1,
			Delete().Table("users").And("id = ?", 1),
		},
	}

	for _, tc := range tcs {
		t.Run("delete", func(t *testing.T) {
			tc.T(t)
		})
	}
}
