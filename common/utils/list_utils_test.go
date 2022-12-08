package utils

import "testing"

func TestListToUnique(t *testing.T) {
	list := []string{"a", "b", "c", "a", "b", "c"}
	uniqueList := ListToUnique(list)

	if len(uniqueList) != 3 {
		t.Errorf("Expected 3, got %d", len(uniqueList))
	}
}

func BenchmarkListToUnique(b *testing.B) {
	list := []string{"a", "b", "c", "a", "b", "c"}
	for i := 0; i < b.N; i++ {
		ListToUnique(list)
	}
}
