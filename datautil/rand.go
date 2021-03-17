package datautil

type Rng64 interface {
	UInt64() uint64
}

type XorShift64Mul struct {
	a uint64
}

func NewXorShift64Mul(seed uint64) Rng64 {
	if seed == 0 {
		panic("seed must not be 0")
	}
	return &XorShift64Mul{a: seed}
}

func (rng *XorShift64Mul) UInt64() uint64 {
	x := rng.a
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	rng.a = x
	return x * 0x2545F4914F6CDD1D
}

type RandomByteStream struct {
	source Rng64
	buffer uint64
	count  int32
}

func NewRandomByteStream(source Rng64) *RandomByteStream {
	return &RandomByteStream{
		source: source,
		buffer: source.UInt64(),
		count:  8}
}

func (rng *RandomByteStream) Read(p []byte) (n int, err error) {
	p = p[0:0]
	return len(p), nil
}
