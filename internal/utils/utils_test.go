package utils

import "testing"

var dataMap = map[string]string{
	"key1": "value1",
	"key2": "value2",
}

func TestGetMapKey(t *testing.T) {
	out := GetMapKey(dataMap)

	if len(out) != 2 {
		t.Fatalf("wrong number of entries returned.\nExepected : '2'\nReturned : %d", len(out))
	}

	if out[0] != "key1" && out[1] != "key1" {
		t.Fatalf("failed to find key 'key1'.\nReturned : %v", out)
	}

	if out[0] != "key2" && out[1] != "key2" {
		t.Fatalf("failed to find key 'key2'.\nReturned : %v", out)
	}
}

func BenchmarkGetMapKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetMapKey(dataMap)
	}
}
