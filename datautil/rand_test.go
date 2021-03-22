package datautil_test

import (
	"io"
	"os"
	"testing"

	"github.com/soerenkoehler/go-testutil/datautil"
)

func TestByteSTream(t *testing.T) {
	in := datautil.NewRng64ByteStream(datautil.NewXorShift64Mul(1))
	if out, err := os.Create("output.random.bin"); err == nil {
		defer out.Close()
		io.CopyN(out, in, 0x10000)
	}
}
