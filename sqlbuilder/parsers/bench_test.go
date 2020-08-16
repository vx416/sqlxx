package parsers

import (
	"testing"
	"time"
)

var benchData = source{ID: 1, Name: "test", CreatedAt: time.Now()}

func BenchmarkParse(b *testing.B) {

	for i := 0; i < b.N; i++ {
		if _, err := New(benchData, false); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQueryParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := make(map[string]interface{})
		if err := ParseConditions(benchData, result, false); err != nil {
			b.Fatal(err)
		}
	}
}
