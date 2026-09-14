package erasure

import (
	"bytes"
	"math/rand"
	"testing"
)

// TestGeneratorIsPrimitive proves the chosen generator walks every nonzero field
// element. Reed Solomon reconstruction breaks silently if it does not.
func TestGeneratorIsPrimitive(t *testing.T) {
	t.Parallel()
	seen := make(map[byte]int)
	for power := 0; power < fieldElements-1; power++ {
		seen[expTable[power]]++
	}
	if len(seen) != fieldElements-1 {
		t.Fatalf("exponent table covers %d distinct values, want %d", len(seen), fieldElements-1)
	}
	if _, ok := seen[0]; ok {
		t.Fatal("exponent table contains zero")
	}
	for value, count := range seen {
		if count != 1 {
			t.Fatalf("value %d appears %d times", value, count)
		}
	}
	if fieldGenerator != 2 || expTable[1] != fieldGenerator {
		t.Fatalf("expTable[1] = %d, want generator %d", expTable[1], fieldGenerator)
	}
}

func TestFieldArithmeticIdentities(t *testing.T) {
	t.Parallel()
	for left := 1; left < fieldElements; left++ {
		if got := mul(byte(left), inverse(byte(left))); got != 1 {
			t.Fatalf("mul(%d, inverse(%d)) = %d, want 1", left, left, got)
		}
		if got := div(byte(left), byte(left)); got != 1 {
			t.Fatalf("div(%d, %d) = %d, want 1", left, left, got)
		}
		if got := mul(byte(left), 0); got != 0 {
			t.Fatalf("mul(%d, 0) = %d, want 0", left, got)
		}
		for right := 1; right < fieldElements; right++ {
			product := mul(byte(left), byte(right))
			if product == 0 {
				t.Fatalf("mul(%d, %d) = 0 for nonzero operands", left, right)
			}
			if got := div(product, byte(right)); got != byte(left) {
				t.Fatalf("div(mul(%d, %d), %d) = %d", left, right, right, got)
			}
			if mul(byte(left), byte(right)) != mul(byte(right), byte(left)) {
				t.Fatalf("mul is not commutative at %d, %d", left, right)
			}
		}
	}
}

func TestMatrixInvertRoundTrip(t *testing.T) {
	t.Parallel()
	source := vandermondeMatrix(6, 4).subMatrix([]int{1, 2, 4, 5})
	inverted, err := source.invert()
	if err != nil {
		t.Fatalf("invert() error = %v", err)
	}
	product := source.multiply(inverted)
	identity := identityMatrix(4)
	for row := range product {
		if !bytes.Equal(product[row], identity[row]) {
			t.Fatalf("row %d = %v, want %v", row, product[row], identity[row])
		}
	}
}

func TestMatrixInvertRejectsSingular(t *testing.T) {
	t.Parallel()
	singular := matrix{{1, 2}, {2, 4}}
	if _, err := singular.invert(); err == nil {
		t.Fatal("invert() error = nil, want singular matrix rejection")
	}
}

// TestReconstructRecoversEveryRecoverableLossPattern removes every combination
// of up to ParityShards shards and requires exact recovery each time.
func TestReconstructRecoversEveryRecoverableLossPattern(t *testing.T) {
	t.Parallel()
	const dataShards, parityShards, shardSize = 4, 3, 97
	codec, err := New(dataShards, parityShards)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	total := codec.TotalShards()
	random := rand.New(rand.NewSource(7))
	original := make([][]byte, total)
	for index := range original {
		original[index] = make([]byte, shardSize)
	}
	for index := 0; index < dataShards; index++ {
		random.Read(original[index])
	}
	if err := codec.Encode(original); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	for mask := 0; mask < 1<<total; mask++ {
		missing := 0
		for index := 0; index < total; index++ {
			if mask&(1<<index) != 0 {
				missing++
			}
		}
		if missing == 0 || missing > parityShards {
			continue
		}
		shards := make([][]byte, total)
		present := make([]bool, total)
		for index := 0; index < total; index++ {
			shards[index] = make([]byte, shardSize)
			if mask&(1<<index) == 0 {
				copy(shards[index], original[index])
				present[index] = true
			}
		}
		if err := codec.Reconstruct(shards, present); err != nil {
			t.Fatalf("Reconstruct(mask=%b) error = %v", mask, err)
		}
		for index := 0; index < total; index++ {
			if !bytes.Equal(shards[index], original[index]) {
				t.Fatalf("mask=%b shard %d not recovered", mask, index)
			}
		}
	}
}

func TestReconstructRejectsInsufficientShards(t *testing.T) {
	t.Parallel()
	codec, err := New(3, 2)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	shards := make([][]byte, codec.TotalShards())
	present := make([]bool, codec.TotalShards())
	for index := range shards {
		shards[index] = make([]byte, 8)
	}
	present[0] = true
	present[4] = true
	if err := codec.Reconstruct(shards, present); err == nil {
		t.Fatal("Reconstruct() error = nil, want insufficient shards")
	}
}

func TestCodecRejectsInvalidParameters(t *testing.T) {
	t.Parallel()
	if _, err := New(0, 2); err == nil {
		t.Fatal("New(0, 2) error = nil, want rejection")
	}
	if _, err := New(4, 0); err == nil {
		t.Fatal("New(4, 0) error = nil, want rejection")
	}
	if _, err := New(200, 100); err == nil {
		t.Fatal("New(200, 100) error = nil, want shard limit rejection")
	}
	codec, err := New(2, 1)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := codec.Encode([][]byte{make([]byte, 4), make([]byte, 5), make([]byte, 4)}); err == nil {
		t.Fatal("Encode() error = nil, want shard length rejection")
	}
	if err := codec.Encode([][]byte{make([]byte, 4), make([]byte, 4)}); err == nil {
		t.Fatal("Encode() error = nil, want shard count rejection")
	}
}

func TestEncodeIsSystematic(t *testing.T) {
	t.Parallel()
	codec, err := New(3, 2)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	identity := identityMatrix(3)
	for row := 0; row < 3; row++ {
		if !bytes.Equal(codec.encodeMatrix[row], identity[row]) {
			t.Fatalf("encode row %d = %v, want identity row", row, codec.encodeMatrix[row])
		}
	}
}
