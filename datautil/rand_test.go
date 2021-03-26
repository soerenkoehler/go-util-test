package datautil_test

import (
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"testing"

	"github.com/soerenkoehler/go-testutil/datautil"
	"github.com/soerenkoehler/go-testutil/testutil"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func TestXorShiftPanicsWithZeroSeed(t *testing.T) {
	testutil.ShouldPanic(t, func() {
		datautil.NewXorShift64Mul(0)
	})
}

func TestSameSeedYieldsSameSequence(t *testing.T) {
	r1 := datautil.NewXorShift64Mul(1)
	r2 := datautil.NewXorShift64Mul(1)
	for i := 0; i < 0x1000_0000; i++ {
		if r1.UInt64() != r2.UInt64() {
			t.Fatalf("Sequences differ at index %d", i)
		}
	}
}

func TestByteDistribution(t *testing.T) {
	seed := uint64(0xffff_ffff_ffff_ffff)
	sampleSize := 0x1_0000_0000
	outcomeNumber := 0x100

	// init expectaion and observation arrays
	exp := make([]float64, outcomeNumber)
	obs := make([]float64, outcomeNumber)
	for i := 0; i < 256; i++ {
		exp[i] = float64(sampleSize / outcomeNumber)
		obs[i] = 0
	}

	// create observations
	src := datautil.NewRng64ByteStream(datautil.NewXorShift64Mul(seed))
	for i := 0; i < sampleSize; i++ {
		b := src.Byte()
		obs[b]++
	}

	// calculate chi square probability
	sum := stat.ChiSquare(obs, exp)
	prob := distuv.ChiSquared{K: float64(outcomeNumber - 1)}.Survival(sum)
	probLimit := 0.9 // may not work with other seed values
	fmt.Printf(
		"Chi Square Sum = %g Probability=%g (limit=%g)\n",
		sum, prob, probLimit)

	if prob < probLimit {
		t.Errorf("confidence below expected limit")
	}

	// normalize observations for entropy calculation
	for i := range obs {
		obs[i] /= float64(sampleSize)
	}

	ent := stat.Entropy(obs) / math.Log(2) // convert to bits
	entLimit := 7.99999
	fmt.Printf("Entropy = %g bits (limit=%g)\n", ent, entLimit)

	if ent < entLimit {
		t.Errorf("entropy below expected limit")
	}
}

func TestRandomDataFile(t *testing.T) {
	seed := uint64(0xffff_ffff_ffff_ffff)
	sampleSize := 0x4000_0000
	fileName := path.Join(t.TempDir(), "rand.bin")

	datautil.CreateRandomFile(fileName, int64(sampleSize), seed)

	in, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("can't open file %s => %v", fileName, err)
	}
	defer in.Close()

	buf, err := io.ReadAll(in)
	if err != nil {
		t.Fatalf("can't read file => %v", err)
	}

	if len(buf) != sampleSize {
		t.Fatalf("expected %d bytes in file but got %d", sampleSize, len(buf))
	}

	cmp := datautil.NewRng64ByteStream(datautil.NewXorShift64Mul(seed))
	for i := range buf {
		if buf[i] != cmp.Byte() {
			t.Fatalf("file differs an position %d", i)
		}
	}
}
