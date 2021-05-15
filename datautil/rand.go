package datautil

import (
	"io"
	"os"
)

type Rng64 interface {
	UInt64() uint64
}

type xorShift64Mul struct {
	a uint64
}

func NewXorShift64Mul(seed uint64) Rng64 {
	// validate seed
	if seed == 0 {
		panic("seed must not be 0")
	}
	// create RNG
	rng := &xorShift64Mul{a: seed}
	// skip first few results
	for i := 13; i > 0; i-- {
		rng.UInt64()
	}
	return rng
}

func (rng *xorShift64Mul) UInt64() uint64 {
	x := rng.a
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	rng.a = x
	return x * 0x2545F4914F6CDD1D
}

type ByteStream interface {
	Byte() byte
}

type Rng64ByteStream struct {
	source Rng64
	buffer uint64
	count  uint32
}

func NewRng64ByteStream(source Rng64) *Rng64ByteStream {
	return &Rng64ByteStream{
		source: source,
		buffer: 0,
		count:  0}
}

func (bs *Rng64ByteStream) Byte() byte {
	if bs.count == 0 {
		bs.buffer = bs.source.UInt64()
		bs.count = 8
	}
	result := byte(bs.buffer)
	bs.buffer >>= 8
	bs.count--
	return result
}

func (bs *Rng64ByteStream) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = bs.Byte()
	}
	return len(p), nil
}

// CreateRandomFile ... TODO
func CreateRandomFile(path string, size int64, seed uint64) {
	in := NewRng64ByteStream(NewXorShift64Mul(seed))
	if out, err := os.Create(path); err == nil {
		defer out.Close()
		io.CopyN(out, in, size)
	} else {
		panic(err)
	}
}
