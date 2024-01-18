package stack

import "testing"

func TestPush(t *testing.T) {
	for i := 1; i <= 3; i++ {
		Push(i)
	}

	want := []int{1, 2, 3}
	got := GetStack()
	for i := 0; i < len(want); i++ {
		if want[i] != got[i] {
			t.Error("elements don't match")
		}
	}
}

func TestPop(t *testing.T) {
	for i := 1; i <= 3; i++ {
	}
}
