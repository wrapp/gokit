package util

import (
	"errors"
	"testing"
)

func TestShortCircuit(t *testing.T) {
	t.Parallel()
	e1 := errors.New("first error")
	e2 := errors.New("second error")
	err := ShortCircuit(
		func() error { return e1 },
		func() error { return e2 },
	)
	if err != e1 {
		t.Errorf("Expected first error got %q", err)
	}

	err = ShortCircuit(
		func() error { return nil },
		func() error { return e2 },
	)
	if err != e2 {
		t.Errorf("Expected second error got %q", err)
	}
}
