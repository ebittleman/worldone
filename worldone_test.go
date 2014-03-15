package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	expected, actual := 3, Add(1, 2)

	if expected != actual {
		t.Errorf("Expected %d But Got %d", expected, actual)
	}
}
