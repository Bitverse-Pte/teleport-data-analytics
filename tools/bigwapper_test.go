package tools

import "testing"

func TestSum(t *testing.T) {
	a := "1"
	b := "2"
	if Sum(a, b) != "3" {
		t.Error("Sum error")
	}
}
