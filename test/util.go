package test

import "testing"

func ShouldPanic(t *testing.T, f func()) {
	defer func() { recover() }()
	f()
	t.Fatal("should have panicked")
}
