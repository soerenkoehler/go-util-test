package data

import (
	"fmt"
	"io"
	"os"
	"path"
)

// TODO always use testing context
func CreateRandomFile(filepath string, size int64, seed uint64) {
	CreateTempFile(filepath, func(out io.Writer) {
		io.CopyN(out, NewRng64ByteStream(NewXorShift64Mul(seed)), size)
	})
}

func CreateTextFile(filepath, content string) {
	CreateTempFile(filepath, func(out io.Writer) {
		fmt.Fprint(out, content)
	})
}

func CreateTempFile(filepath string, contentWriter func(io.Writer)) {
	err := os.MkdirAll(path.Dir(filepath), 0700)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	defer out.Close()

	contentWriter(out)
}
