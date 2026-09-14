package audit

import (
	"bytes"
	"errors"
	"math"
	"math/big"
	"testing"
)

// exactMissProbability computes C(n-d, k)/C(n, k) with exact rational
// arithmetic, so the floating point implementation is checked against ground
// truth rather than against itself.
func exactMissProbability(chunks, bad, samples int) *big.Rat {
	if bad == 0 {
		return big.NewRat(1, 1)
	}
	if samples > chunks-bad {
		return new(big.Rat)
	}
	numerator := new(big.Int).Binomial(int64(chunks-bad), int64(samples))
	denominator := new(big.Int).Binomial(int64(chunks), int64(samples))
	return new(big.Rat).SetFrac(numerator, denominator)
}

func TestMissProbabilityMatchesExactHypergeometric(t *testing.T) {
	t.Parallel()
	for chunks := 1; chunks <= 20; chunks++ {
		for bad := 0; bad <= chunks; bad++ {
			for samples := 0; samples <= chunks; samples++ {
				got, err := MissProbability(chunks, bad, samples)
				if err != nil {
					t.Fatalf("MissProbability(%d, %d, %d) error = %v", chunks, bad, samples, err)
				}
				want, _ := exactMissProbability(chunks, bad, samples).Float64()
				if math.Abs(got-want) > 1e-12 {
					t.Fatalf("MissProbability(%d, %d, %d) = %.15f, want %.15f", chunks, bad, samples, got, want)
				}
			}
		}
	}
}

func TestMissProbabilityIsMonotonic(t *testing.T) {
	t.Parallel()
	const chunks = 500
	previous := 1.0
	for samples := 0; samples <= 50; samples++ {
		probability, err := MissProbability(chunks, 5, samples)
		if err != nil {
			t.Fatalf("MissProbability() error = %v", err)
		}
		if probability > previous+1e-15 {
			t.Fatalf("probability rose from %v to %v at %d samples", previous, probability, samples)
		}
		previous = probability
	}
	fewer, err := MissProbability(chunks, 1, 20)
	if err != nil {
		t.Fatalf("MissProbability() error = %v", err)
	}
	more, err := MissProbability(chunks, 10, 20)
	if err != nil {
		t.Fatalf("MissProbability() error = %v", err)
	}
	if more >= fewer {
		t.Fatalf("more damage did not lower the miss probability: %v vs %v", more, fewer)
	}
}

func TestRequiredSamplesIsTheSmallestSufficientCount(t *testing.T) {
	t.Parallel()
	cases := []struct {
		chunks     int
		bad        int
		confidence float64
	}{
		{chunks: 100, bad: 1, confidence: 0.95},
		{chunks: 1000, bad: 10, confidence: 0.99},
		{chunks: 13730, bad: 137, confidence: 0.99},
		{chunks: 64, bad: 32, confidence: 0.999},
	}
	for _, item := range cases {
		samples, err := RequiredSamples(item.chunks, item.bad, item.confidence)
		if err != nil {
			t.Fatalf("RequiredSamples(%+v) error = %v", item, err)
		}
		alpha := 1 - item.confidence
		reached, err := MissProbability(item.chunks, item.bad, samples)
		if err != nil {
			t.Fatalf("MissProbability() error = %v", err)
		}
		if reached > alpha {
			t.Fatalf("RequiredSamples(%+v) = %d leaves probability %v above %v", item, samples, reached, alpha)
		}
		if samples > 1 {
			short, err := MissProbability(item.chunks, item.bad, samples-1)
			if err != nil {
				t.Fatalf("MissProbability() error = %v", err)
			}
			if short <= alpha {
				t.Fatalf("RequiredSamples(%+v) = %d is not minimal", item, samples)
			}
		}
	}
}

func TestSamplesForFractionMatchesDocumentedExample(t *testing.T) {
	t.Parallel()
	// One percent damage detected with 95 percent confidence needs roughly 300
	// samples, and 99 percent needs roughly 460, for a large object.
	samples, err := SamplesForFraction(100000, 0.01, 0.95)
	if err != nil {
		t.Fatalf("SamplesForFraction() error = %v", err)
	}
	if samples < 280 || samples > 320 {
		t.Fatalf("95 percent confidence needs %d samples, want roughly 300", samples)
	}
	samples, err = SamplesForFraction(100000, 0.01, 0.99)
	if err != nil {
		t.Fatalf("SamplesForFraction() error = %v", err)
	}
	if samples < 440 || samples > 480 {
		t.Fatalf("99 percent confidence needs %d samples, want roughly 460", samples)
	}
}

func TestMaxUndetectedBadChunksBoundsAreConsistent(t *testing.T) {
	t.Parallel()
	const chunks, confidence = 5000, 0.99
	for _, samples := range []int{0, 1, 10, 100, 460, 1000, chunks} {
		bound, err := MaxUndetectedBadChunks(chunks, samples, confidence)
		if err != nil {
			t.Fatalf("MaxUndetectedBadChunks() error = %v", err)
		}
		if samples == 0 && bound != chunks {
			t.Fatalf("zero samples gave bound %d, want %d", bound, chunks)
		}
		if samples == chunks && bound != 0 {
			t.Fatalf("exhaustive audit gave bound %d, want 0", bound)
		}
		if bound < 0 || bound > chunks {
			t.Fatalf("bound %d is out of range", bound)
		}
		if bound > 0 && bound < chunks {
			// The bound must be the largest damage level still plausible.
			plausible, err := MissProbability(chunks, bound, samples)
			if err != nil {
				t.Fatalf("MissProbability() error = %v", err)
			}
			if plausible <= 1-confidence {
				t.Fatalf("bound %d at %d samples should have been rejected", bound, samples)
			}
			rejected, err := MissProbability(chunks, bound+1, samples)
			if err != nil {
				t.Fatalf("MissProbability() error = %v", err)
			}
			if rejected > 1-confidence {
				t.Fatalf("bound %d at %d samples is not maximal", bound, samples)
			}
		}
	}
	if got := MinIntactFraction(1000, 10); math.Abs(got-0.99) > 1e-12 {
		t.Fatalf("MinIntactFraction(1000, 10) = %v, want 0.99", got)
	}
	if got := MinIntactFraction(1000, 1000); got != 0 {
		t.Fatalf("MinIntactFraction with total damage = %v, want 0", got)
	}
}

func TestSelectSamplesIsDistinctInRangeAndSorted(t *testing.T) {
	t.Parallel()
	indexes, err := SelectSamples(1000, 64, nil)
	if err != nil {
		t.Fatalf("SelectSamples() error = %v", err)
	}
	if len(indexes) != 64 {
		t.Fatalf("len(indexes) = %d, want 64", len(indexes))
	}
	seen := make(map[int]struct{}, len(indexes))
	previous := -1
	for _, index := range indexes {
		if index < 0 || index >= 1000 {
			t.Fatalf("index %d out of range", index)
		}
		if _, exists := seen[index]; exists {
			t.Fatalf("index %d repeated", index)
		}
		seen[index] = struct{}{}
		if index <= previous {
			t.Fatalf("indexes are not sorted at %d", index)
		}
		previous = index
	}

	all, err := SelectSamples(8, 20, nil)
	if err != nil {
		t.Fatalf("SelectSamples() error = %v", err)
	}
	if len(all) != 8 {
		t.Fatalf("oversized request returned %d indexes, want 8", len(all))
	}
}

func TestSelectSamplesVariesBetweenAudits(t *testing.T) {
	t.Parallel()
	first, err := SelectSamples(100000, 32, nil)
	if err != nil {
		t.Fatalf("SelectSamples() error = %v", err)
	}
	second, err := SelectSamples(100000, 32, nil)
	if err != nil {
		t.Fatalf("SelectSamples() error = %v", err)
	}
	identical := len(first) == len(second)
	if identical {
		for index := range first {
			if first[index] != second[index] {
				identical = false
				break
			}
		}
	}
	if identical {
		t.Fatal("two audits drew identical samples, which a holder could predict")
	}
}

func TestSelectSamplesReportsSourceFailure(t *testing.T) {
	t.Parallel()
	if _, err := SelectSamples(100, 4, failingReader{}); err == nil {
		t.Fatal("SelectSamples() error = nil, want randomness failure")
	}
	if _, err := SelectSamples(100, 4, bytes.NewReader(nil)); err == nil {
		t.Fatal("SelectSamples() error = nil, want exhausted source failure")
	}
}

func TestRejectsInvalidParameters(t *testing.T) {
	t.Parallel()
	if _, err := MissProbability(0, 0, 0); err == nil {
		t.Fatal("MissProbability(0, ...) error = nil")
	}
	if _, err := MissProbability(10, 11, 1); err == nil {
		t.Fatal("MissProbability with excess damage error = nil")
	}
	if _, err := MissProbability(10, 1, 11); err == nil {
		t.Fatal("MissProbability with excess samples error = nil")
	}
	if _, err := RequiredSamples(10, 1, 1); err == nil {
		t.Fatal("RequiredSamples with confidence 1 error = nil")
	}
	if _, err := RequiredSamples(10, 0, 0.99); err == nil {
		t.Fatal("RequiredSamples with zero damage error = nil")
	}
	if _, err := SamplesForFraction(10, 0, 0.99); err == nil {
		t.Fatal("SamplesForFraction with zero tolerance error = nil")
	}
	if _, err := MaxUndetectedBadChunks(10, 11, 0.99); err == nil {
		t.Fatal("MaxUndetectedBadChunks with excess samples error = nil")
	}
	if _, err := SelectSamples(0, 1, nil); err == nil {
		t.Fatal("SelectSamples with zero chunks error = nil")
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errors.New("no randomness available")
}
