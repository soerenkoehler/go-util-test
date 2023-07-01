package data_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"testing"

	"github.com/soerenkoehler/go-util-test/data"
	"github.com/soerenkoehler/go-util-test/test"
)

func TestXorShiftPanicsWithZeroSeed(t *testing.T) {
	test.ShouldPanic(t, func() {
		data.NewXorShift64Mul(0)
	})
}

func TestSameSeedYieldsSameSequence(t *testing.T) {
	r1 := data.NewXorShift64Mul(1)
	r2 := data.NewXorShift64Mul(1)
	for i := 0; i < 0x1000_0000; i++ {
		if r1.UInt64() != r2.UInt64() {
			t.Fatalf("Sequences differ at index %d", i)
		}
	}
}

type RegressionTestCase struct {
	seed uint64
	size int64
	hash string
}

// TestSampleRegressions does not evaluate the quality of the RNG but does a
// regression test over some data samples to safeguard against unintended
// modifications of the RNG code.
func TestSampleRegressions(t *testing.T) {
	samples := []RegressionTestCase{
		{0x0000_0000_0000_0001, 0x10_0000, "9c236b00023d7e109be93c569223c1ae7104465ca18909bb5a72e981e42b746b"},
		{0x0000_0000_0000_0020, 0x10_0000, "b2f62151f1d9b18e053f3a401cb1caf71dd08057bf93a8df94bb1f2961a408e8"},
		{0x0000_0000_0000_0400, 0x10_0000, "ebf3066836840b75f9f308a9dc7f2ee4a97031f736d192b0727cecab7da82f1c"},
		{0x0000_0000_0000_8000, 0x10_0000, "4496d27fdc6d61a4e04a3dccf1066f1c926eef997cf287776936b3bdb4d82efb"},

		{0x0000_0000_0001_0000, 0x10_0000, "88cacdf18e57df498b1fd59207ce3bac79f89a69a1dbc40d8e318fab1afb87dd"},
		{0x0000_0000_0020_0000, 0x10_0000, "07fbc39374395094119c3fb7caff38451d94ff8ea046d5821732059c09fb985e"},
		{0x0000_0000_0400_0000, 0x10_0000, "ec1fdd5c5352a89d20912962b386890fc71efa6ccd528f8e66e8fac392a5bf5c"},
		{0x0000_0000_8000_0000, 0x10_0000, "c35ac2ed2abdf1982a864ff38cf3c69ed75fbf6d1adb33b90fe310c575b4bd8b"},

		{0x0000_0001_0000_0000, 0x10_0000, "f92e7d0b695bd5b4479dc3c65f82c968dd1e87f699abba85a618b0f279aec101"},
		{0x0000_0020_0000_0000, 0x10_0000, "ccc8d94de64591f7e975bfcc301d74a1a38a86d4b3dc7e7c016bd9f99130aec7"},
		{0x0000_0400_0000_0000, 0x10_0000, "2f0a1e0fd48aa0b434ff0566f30d7fd74affde62dd9bb8a07e288dc4ffd18599"},
		{0x0000_8000_0000_0000, 0x10_0000, "ca66c20d634a50f6dd9434ffb81cdb9489da37da026dd8b243dea3d0b2dc74cc"},

		{0x0001_0000_0000_0000, 0x10_0000, "cd746664b52eec7ebac1891259f7c879a3ab7d1e2d1da5d767ec9067e5c646cd"},
		{0x0020_0000_0000_0000, 0x10_0000, "e631135c4e69e201286403038da42950dc0ab40b96b144907a07b0aca1bfefdb"},
		{0x0400_0000_0000_0000, 0x10_0000, "8ae6d104e194022cc616bbae94f8493367b98a26acb781e5ac3a86b13508963a"},
		{0x8000_0000_0000_0000, 0x10_0000, "a27a1b8ab09f4ff66992a45e05713c78d274b1263ce4d76e201d0c506008fa3c"}}

	for i, sample := range samples {
		rng := data.NewRng64ByteStream(data.NewXorShift64Mul(sample.seed))
		hash := sha256.New()

		io.CopyN(hash, rng, sample.size)
		hashValue := hex.EncodeToString(hash.Sum(nil))

		if hashValue != sample.hash {
			t.Fatalf(
				"Mismatch in sample %d:\nexpected: %s\nactual: %s\n",
				i,
				sample.hash,
				hashValue)
		}
	}
}

func TestRandomDataFile(t *testing.T) {
	seed := uint64(0xffff_ffff_ffff_ffff)
	sampleSize := 0x4000_0000
	fileName := path.Join(t.TempDir(), "rand.bin")

	data.CreateRandomFile(fileName, int64(sampleSize), seed)

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

	cmp := data.NewRng64ByteStream(data.NewXorShift64Mul(seed))
	for i := range buf {
		if buf[i] != cmp.Byte() {
			t.Fatalf("file differs on position %d", i)
		}
	}
}
