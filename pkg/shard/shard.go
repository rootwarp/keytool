// Package shard provides data splitting and reconstruction using erasure coding.
package shard

// Splitter defines the interface for splitting and reconstructing data.
type Splitter interface {
	// Split divides data into multiple shards (data + parity).
	// Returns shards that can tolerate loss of up to parityShards pieces.
	Split(data []byte) ([][]byte, error)

	// Reconstruct recovers original data from shards.
	// Some shards may be nil (missing), as long as enough remain for recovery.
	Reconstruct(shards [][]byte) ([]byte, error)

	// TotalShards returns the total number of shards (data + parity).
	TotalShards() int
}

// Store defines the interface for shard persistence.
type Store interface {
	// Save persists shards for a given account.
	// Files are named as <account>.<shard_index>.
	Save(account string, shards [][]byte) error

	// Load retrieves shards for a given account.
	// totalShards specifies how many shards to expect.
	// Missing shards are returned as nil in the slice.
	Load(account string, totalShards int) ([][]byte, error)
}
