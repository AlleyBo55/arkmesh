package erasure

import (
	"errors"
	"fmt"
)

const MaxShards = 255

// Codec is a systematic Reed Solomon coder: the first DataShards shards are the
// original data unchanged, and any DataShards of the total shards reconstruct
// everything.
type Codec struct {
	DataShards   int
	ParityShards int
	encodeMatrix matrix
}

func New(dataShards, parityShards int) (*Codec, error) {
	if dataShards < 1 {
		return nil, errors.New("erasure: at least one data shard is required")
	}
	if parityShards < 1 {
		return nil, errors.New("erasure: at least one parity shard is required")
	}
	if dataShards+parityShards > MaxShards {
		return nil, fmt.Errorf("erasure: %d shards exceed the %d shard limit", dataShards+parityShards, MaxShards)
	}
	total := dataShards + parityShards
	vandermonde := vandermondeMatrix(total, dataShards)
	top := make([]int, dataShards)
	for index := range top {
		top[index] = index
	}
	inverted, err := vandermonde.subMatrix(top).invert()
	if err != nil {
		return nil, err
	}
	return &Codec{
		DataShards:   dataShards,
		ParityShards: parityShards,
		encodeMatrix: vandermonde.multiply(inverted),
	}, nil
}

func (codec *Codec) TotalShards() int {
	return codec.DataShards + codec.ParityShards
}

// Encode fills the parity shards from the data shards. Every shard must already
// be allocated with the same length.
func (codec *Codec) Encode(shards [][]byte) error {
	shardSize, err := codec.checkShards(shards)
	if err != nil {
		return err
	}
	for parity := 0; parity < codec.ParityShards; parity++ {
		row := codec.encodeMatrix[codec.DataShards+parity]
		target := shards[codec.DataShards+parity]
		for index := range target {
			target[index] = 0
		}
		for data := 0; data < codec.DataShards; data++ {
			factor := row[data]
			if factor == 0 {
				continue
			}
			source := shards[data]
			for offset := 0; offset < shardSize; offset++ {
				target[offset] ^= mul(factor, source[offset])
			}
		}
	}
	return nil
}

// Reconstruct rebuilds every missing shard from any DataShards present shards.
// Shards marked absent are overwritten; present shards are never modified.
func (codec *Codec) Reconstruct(shards [][]byte, present []bool) error {
	if len(present) != codec.TotalShards() {
		return fmt.Errorf("erasure: expected %d presence flags, got %d", codec.TotalShards(), len(present))
	}
	shardSize, err := codec.checkShards(shards)
	if err != nil {
		return err
	}
	available := make([]int, 0, codec.DataShards)
	for index, ok := range present {
		if ok && len(available) < codec.DataShards {
			available = append(available, index)
		}
	}
	if len(available) < codec.DataShards {
		return fmt.Errorf("erasure: %d shards present, %d required", len(available), codec.DataShards)
	}

	decodeMatrix, err := codec.encodeMatrix.subMatrix(available).invert()
	if err != nil {
		return err
	}
	recovered := make([][]byte, codec.DataShards)
	for data := 0; data < codec.DataShards; data++ {
		if present[data] {
			recovered[data] = shards[data]
			continue
		}
		target := make([]byte, shardSize)
		for position, source := range available {
			factor := decodeMatrix[data][position]
			if factor == 0 {
				continue
			}
			shard := shards[source]
			for offset := 0; offset < shardSize; offset++ {
				target[offset] ^= mul(factor, shard[offset])
			}
		}
		recovered[data] = target
	}
	for data := 0; data < codec.DataShards; data++ {
		if !present[data] {
			copy(shards[data], recovered[data])
			present[data] = true
		}
	}
	for parity := 0; parity < codec.ParityShards; parity++ {
		index := codec.DataShards + parity
		if present[index] {
			continue
		}
		row := codec.encodeMatrix[index]
		target := shards[index]
		for offset := range target {
			target[offset] = 0
		}
		for data := 0; data < codec.DataShards; data++ {
			factor := row[data]
			if factor == 0 {
				continue
			}
			source := shards[data]
			for offset := 0; offset < shardSize; offset++ {
				target[offset] ^= mul(factor, source[offset])
			}
		}
		present[index] = true
	}
	return nil
}

func (codec *Codec) checkShards(shards [][]byte) (int, error) {
	if len(shards) != codec.TotalShards() {
		return 0, fmt.Errorf("erasure: expected %d shards, got %d", codec.TotalShards(), len(shards))
	}
	shardSize := len(shards[0])
	if shardSize == 0 {
		return 0, errors.New("erasure: shard size must be positive")
	}
	for _, shard := range shards {
		if len(shard) != shardSize {
			return 0, errors.New("erasure: all shards must have the same length")
		}
	}
	return shardSize, nil
}
