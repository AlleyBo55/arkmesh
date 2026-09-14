// Package audit computes the statistics behind sampled retrievability checks.
//
// A holder that keeps only the chunks it was asked for last time would pass a
// predictable audit, so samples are drawn from a cryptographic source. The
// bounds below come from sampling WITHOUT replacement, which is the
// hypergeometric distribution, not the looser binomial approximation:
//
//	P(all k sampled chunks are clean | d of n chunks are bad)
//	  = C(n-d, k) / C(n, k)
//	  = product over i in [0, k) of (n-d-i) / (n-i)
//
// That probability is what lets a clean audit state an upper bound on how much
// of an object can still be missing.
package audit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"sort"
)

// MissProbability returns the probability that k chunks sampled without
// replacement all come back clean while d of n chunks are actually bad.
func MissProbability(chunks, bad, samples int) (float64, error) {
	if chunks <= 0 {
		return 0, errors.New("audit: chunk count must be positive")
	}
	if bad < 0 || bad > chunks {
		return 0, fmt.Errorf("audit: bad chunk count %d is out of range", bad)
	}
	if samples < 0 || samples > chunks {
		return 0, fmt.Errorf("audit: sample count %d is out of range", samples)
	}
	if bad == 0 {
		return 1, nil
	}
	if samples > chunks-bad {
		// More samples than clean chunks: at least one bad chunk is certain.
		return 0, nil
	}
	logProbability := 0.0
	for index := 0; index < samples; index++ {
		logProbability += math.Log(float64(chunks-bad-index)) - math.Log(float64(chunks-index))
	}
	return math.Exp(logProbability), nil
}

// RequiredSamples returns the smallest number of samples whose clean result
// would rule out `bad` or more damaged chunks at the given confidence.
func RequiredSamples(chunks, bad int, confidence float64) (int, error) {
	if err := checkConfidence(confidence); err != nil {
		return 0, err
	}
	if bad <= 0 {
		return 0, errors.New("audit: tolerated damage must be at least one chunk")
	}
	alpha := 1 - confidence
	for samples := 1; samples <= chunks; samples++ {
		probability, err := MissProbability(chunks, bad, samples)
		if err != nil {
			return 0, err
		}
		if probability <= alpha {
			return samples, nil
		}
	}
	return chunks, nil
}

// SamplesForFraction converts a tolerated damage fraction into a sample count.
func SamplesForFraction(chunks int, fraction, confidence float64) (int, error) {
	if chunks <= 0 {
		return 0, errors.New("audit: chunk count must be positive")
	}
	if fraction <= 0 || fraction > 1 {
		return 0, fmt.Errorf("audit: tolerance %v must be in (0, 1]", fraction)
	}
	bad := int(math.Ceil(fraction * float64(chunks)))
	if bad < 1 {
		bad = 1
	}
	return RequiredSamples(chunks, bad, confidence)
}

// MaxUndetectedBadChunks reports how many bad chunks could still hide after a
// fully clean audit of `samples` chunks, at the given confidence.
//
// It is the honest reading of a passed audit: not "the object is intact", but
// "no more than this much of it can be missing, with this confidence".
func MaxUndetectedBadChunks(chunks, samples int, confidence float64) (int, error) {
	if err := checkConfidence(confidence); err != nil {
		return 0, err
	}
	if chunks <= 0 {
		return 0, errors.New("audit: chunk count must be positive")
	}
	if samples < 0 || samples > chunks {
		return 0, fmt.Errorf("audit: sample count %d is out of range", samples)
	}
	alpha := 1 - confidence
	for bad := 1; bad <= chunks; bad++ {
		probability, err := MissProbability(chunks, bad, samples)
		if err != nil {
			return 0, err
		}
		if probability <= alpha {
			return bad - 1, nil
		}
	}
	return chunks, nil
}

// MinIntactFraction converts an undetected damage bound into a lower bound on
// how much of the object is present.
func MinIntactFraction(chunks, maxBad int) float64 {
	if chunks <= 0 {
		return 0
	}
	if maxBad >= chunks {
		return 0
	}
	return float64(chunks-maxBad) / float64(chunks)
}

// SelectSamples draws `samples` distinct chunk indexes without replacement from
// a cryptographic source. Predictable selection would let a holder store only
// the chunks it expects to be asked about, so this never uses a seeded PRNG.
func SelectSamples(chunks, samples int, random io.Reader) ([]int, error) {
	if chunks <= 0 {
		return nil, errors.New("audit: chunk count must be positive")
	}
	if samples < 0 {
		return nil, errors.New("audit: sample count must not be negative")
	}
	if random == nil {
		random = rand.Reader
	}
	if samples >= chunks {
		all := make([]int, chunks)
		for index := range all {
			all[index] = index
		}
		return all, nil
	}
	limit := big.NewInt(int64(chunks))
	chosen := make(map[int]struct{}, samples)
	for len(chosen) < samples {
		drawn, err := rand.Int(random, limit)
		if err != nil {
			return nil, fmt.Errorf("audit: draw sample: %w", err)
		}
		chosen[int(drawn.Int64())] = struct{}{}
	}
	indexes := make([]int, 0, samples)
	for index := range chosen {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	return indexes, nil
}

func checkConfidence(confidence float64) error {
	if confidence <= 0 || confidence >= 1 {
		return fmt.Errorf("audit: confidence %v must be in (0, 1)", confidence)
	}
	return nil
}
