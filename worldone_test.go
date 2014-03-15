package worldone

import (
	"testing"
)

func TestAdd(t *testing.T) {
	expected, actual := 2, Add(1, 1)

	if expected != actual {
		t.Errorf("Expected %d But Got %d", expected, actual)
	}
}
