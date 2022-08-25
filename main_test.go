package main

import "testing"

func TestXxx(t *testing.T) {
	if 1 != 2 {
		t.Errorf("1 is not equal to 2")
	}
}
