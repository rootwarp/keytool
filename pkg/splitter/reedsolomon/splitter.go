package reedsolomon

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/klauspost/reedsolomon"
	"github.com/rootwarp/keytool/pkg/shard"
)

// reedSolomonSplitter implements shard.Splitter using Reed-Solomon erasure coding.
type reedSolomonSplitter struct {
	encoder      reedsolomon.Encoder
	dataShards   int
	parityShards int
}

// NewSplitter creates a new Splitter with the specified parameters.
// dataShards is the number of data shards (k).
// parityShards is the number of parity shards (m).
func NewSplitter(dataShards, parityShards int) (shard.Splitter, error) {
	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}

	return &reedSolomonSplitter{
		encoder:      enc,
		dataShards:   dataShards,
		parityShards: parityShards,
	}, nil
}

// Split divides data into multiple shards using Reed-Solomon encoding.
func (s *reedSolomonSplitter) Split(data []byte) ([][]byte, error) {
	// Prepend original size to data for reconstruction
	sizeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sizeBytes, uint64(len(data)))
	dataWithSize := append(sizeBytes, data...)

	// Split data into shards
	shards, err := s.encoder.Split(dataWithSize)
	if err != nil {
		return nil, fmt.Errorf("failed to split data: %w", err)
	}

	// Encode parity shards
	if err := s.encoder.Encode(shards); err != nil {
		return nil, fmt.Errorf("failed to encode shards: %w", err)
	}

	return shards, nil
}

// Reconstruct recovers original data from shards.
func (s *reedSolomonSplitter) Reconstruct(shards [][]byte) ([]byte, error) {
	// Reconstruct missing shards
	if err := s.encoder.Reconstruct(shards); err != nil {
		return nil, fmt.Errorf("failed to reconstruct shards: %w", err)
	}

	// Determine shard size from available shards
	var shardSize int
	for _, shard := range shards {
		if shard != nil {
			shardSize = len(shard)
			break
		}
	}

	// Join shards back together
	totalSize := shardSize * s.dataShards
	var buf bytes.Buffer
	if err := s.encoder.Join(&buf, shards, totalSize); err != nil {
		return nil, fmt.Errorf("failed to join shards: %w", err)
	}

	// Extract original size from header
	joined := buf.Bytes()
	if len(joined) < 8 {
		return nil, fmt.Errorf("invalid data: too short")
	}
	originalSize := binary.BigEndian.Uint64(joined[:8])

	// Validate size
	if int(originalSize) > len(joined)-8 {
		return nil, fmt.Errorf("invalid data: size header exceeds available data")
	}

	return joined[8 : 8+originalSize], nil
}

// TotalShards returns the total number of shards (data + parity).
func (s *reedSolomonSplitter) TotalShards() int {
	return s.dataShards + s.parityShards
}
