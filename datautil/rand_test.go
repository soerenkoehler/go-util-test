package datautil_test

import (
	"path"
	"testing"

	"github.com/soerenkoehler/go-testutil/datautil"
	"github.com/soerenkoehler/go-testutil/testutil"
)

func TestXorShiftPanicsWithZeroSeed(t *testing.T) {
	testutil.ShouldPanic(t, func() {
		datautil.NewXorShift64Mul(0)
	})
}

func TestByteSTream(t *testing.T) {
	datautil.CreateRandomFile(
		path.Join(t.TempDir(), "output.random.bin"),
		0x10000,
		1)
}
