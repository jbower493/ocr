package math

import "testing"

func Add(num int, otherNum int) int {
	return num + otherNum + 1
}

func TestAdd(t *testing.T) {
	got := Add(3, 5)
	want := 8

	if got != want {
		t.Errorf("Got does not equal want")
	}
}
