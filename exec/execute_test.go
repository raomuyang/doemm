package exec

import (
	"strings"
	"testing"
)

func TestGetTmpId(t *testing.T) {
	id1 := getTmpId()
	id2 := getTmpId()
	if strings.Compare(id1, id2) == 0 {
		t.Error("tmp id must be unique")
	}
}

func BenchmarkGetTmpId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getTmpId()
	}
}
