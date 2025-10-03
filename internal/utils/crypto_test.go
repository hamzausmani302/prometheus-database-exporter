package utils

import (
	"testing"
)

func TestHashingFunction(t *testing.T) {
	result := Hash("test-metric", "select * from test_table", "test_label1test_label2")
	result2 := Hash("test-metric", "select * from test_table", "test_label1test_label2")
	if result != result2{
		t.Errorf("Expected = %s, got = %s", result, result2)
	}
	result3 := Hash("test-metric", "select * from test_table", "test_label1test_label3")
	if result == result3{
		t.Error("expected the hashes to be different but are same")
	}
}
