package splitter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/rootwarp/keytool/pkg/filestore"
	"github.com/rootwarp/keytool/pkg/shard"
	"github.com/rootwarp/keytool/pkg/splitter/reedsolomon"
)

// Splitter handles splitting and reconstructing secret data.
type Splitter interface {
	// Split encodes the data as base64 JSON and splits it into shards.
	Split(data SecretData, shardDir string) error

	// Reconstruct recovers the original data from shards.
	Reconstruct(address string, shardDir string) (*SecretData, error)

	// TotalShards returns the total number of shards (data + parity).
	TotalShards() int
}

// defaultSplitter implements Splitter using Reed-Solomon encoding.
type defaultSplitter struct {
	rsSplitter shard.Splitter
}

// NewSplitter creates a new Splitter with the specified Reed-Solomon parameters.
func NewSplitter(dataShards, parityShards int) (Splitter, error) {
	rs, err := reedsolomon.NewSplitter(dataShards, parityShards)
	if err != nil {
		return nil, fmt.Errorf("failed to create RS splitter: %w", err)
	}

	return &defaultSplitter{
		rsSplitter: rs,
	}, nil
}

// Split encodes data as JSON, base64 encodes it, splits using RS, and saves to files.
func (s *defaultSplitter) Split(data SecretData, shardDir string) error {
	// Step 1: Marshal to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Step 2: Base64 encode
	encoded := base64.StdEncoding.EncodeToString(jsonData)

	// Step 3: Split using Reed-Solomon
	shards, err := s.rsSplitter.Split([]byte(encoded))
	if err != nil {
		return fmt.Errorf("failed to split data: %w", err)
	}

	// Step 4: Save to filesystem
	store := filestore.NewFileStore(shardDir)
	if err := store.Save(data.Address, shards); err != nil {
		return fmt.Errorf("failed to save shards: %w", err)
	}

	return nil
}

// Reconstruct loads shards, reconstructs the data, and decodes from base64 JSON.
func (s *defaultSplitter) Reconstruct(address string, shardDir string) (*SecretData, error) {
	// Step 1: Load shards from filesystem
	store := filestore.NewFileStore(shardDir)
	shards, err := store.Load(address, s.rsSplitter.TotalShards())
	if err != nil {
		return nil, fmt.Errorf("failed to load shards: %w", err)
	}

	// Step 2: Reconstruct using Reed-Solomon
	encoded, err := s.rsSplitter.Reconstruct(shards)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct data: %w", err)
	}

	// Step 3: Base64 decode
	jsonData, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Step 4: Unmarshal JSON
	var data SecretData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &data, nil
}

// TotalShards returns the total number of shards.
func (s *defaultSplitter) TotalShards() int {
	return s.rsSplitter.TotalShards()
}
