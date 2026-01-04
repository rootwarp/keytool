// Package filestore provides file-based shard persistence.
package filestore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rootwarp/keytool/pkg/shard"
)

// fileStore implements shard.Store using the local filesystem.
type fileStore struct {
	baseDir string
}

// NewFileStore creates a new FileStore with the specified base directory.
func NewFileStore(baseDir string) shard.Store {
	return &fileStore{baseDir: baseDir}
}

// Save persists shards to files named <account>.<shard_index>.
func (s *fileStore) Save(account string, shards [][]byte) error {
	if err := os.MkdirAll(s.baseDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	for i, shard := range shards {
		filename := fmt.Sprintf("%s.%d", account, i)
		filePath := filepath.Join(s.baseDir, filename)

		if err := os.WriteFile(filePath, shard, 0600); err != nil {
			return fmt.Errorf("failed to write shard %d: %w", i, err)
		}
	}

	return nil
}

// Load retrieves shards from files. Missing shards are returned as nil.
func (s *fileStore) Load(account string, totalShards int) ([][]byte, error) {
	shards := make([][]byte, totalShards)

	for i := range totalShards {
		filename := fmt.Sprintf("%s.%d", account, i)
		filePath := filepath.Join(s.baseDir, filename)

		data, err := os.ReadFile(filePath)
		if os.IsNotExist(err) {
			shards[i] = nil
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read shard %d: %w", i, err)
		}

		shards[i] = data
	}

	return shards, nil
}
