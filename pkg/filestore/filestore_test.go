package filestore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rootwarp/keytool/pkg/splitter/reedsolomon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileStore_CreatesInstance(t *testing.T) {
	store := NewFileStore("/tmp/test")

	assert.NotNil(t, store)
}

func TestFileStore_Save_CreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileStore(tmpDir)
	shards := [][]byte{
		[]byte("shard0"),
		[]byte("shard1"),
		[]byte("shard2"),
		[]byte("shard3"),
	}

	err := store.Save("testaccount", shards)

	require.NoError(t, err)
	for i := range shards {
		filename := filepath.Join(tmpDir, "testaccount."+string(rune('0'+i)))
		assert.FileExists(t, filename)
	}
}

func TestFileStore_Save_CreatesDirectoryIfNotExists(t *testing.T) {
	tmpDir := filepath.Join(t.TempDir(), "nonexistent", "path")
	store := NewFileStore(tmpDir)
	shards := [][]byte{[]byte("shard0"), []byte("shard1")}

	err := store.Save("account", shards)

	require.NoError(t, err)
	assert.DirExists(t, tmpDir)
}

func TestFileStore_Save_WritesCorrectContent(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileStore(tmpDir)
	shards := [][]byte{
		[]byte("content0"),
		[]byte("content1"),
	}

	err := store.Save("myaccount", shards)
	require.NoError(t, err)

	content0, err := os.ReadFile(filepath.Join(tmpDir, "myaccount.0"))
	require.NoError(t, err)
	assert.Equal(t, []byte("content0"), content0)

	content1, err := os.ReadFile(filepath.Join(tmpDir, "myaccount.1"))
	require.NoError(t, err)
	assert.Equal(t, []byte("content1"), content1)
}

func TestFileStore_Load_ReturnsShards(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileStore(tmpDir)
	originalShards := [][]byte{
		[]byte("shard0"),
		[]byte("shard1"),
		[]byte("shard2"),
		[]byte("shard3"),
	}
	err := store.Save("account", originalShards)
	require.NoError(t, err)

	loadedShards, err := store.Load("account", 4)

	require.NoError(t, err)
	assert.Equal(t, originalShards, loadedShards)
}

func TestFileStore_Load_ReturnsNilForMissingShards(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileStore(tmpDir)
	// Only write shards 0 and 2
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "account.0"), []byte("shard0"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "account.2"), []byte("shard2"), 0600))

	loadedShards, err := store.Load("account", 4)

	require.NoError(t, err)
	assert.Equal(t, []byte("shard0"), loadedShards[0])
	assert.Nil(t, loadedShards[1])
	assert.Equal(t, []byte("shard2"), loadedShards[2])
	assert.Nil(t, loadedShards[3])
}

func TestFileStore_SaveAndLoad_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileStore(tmpDir)
	originalShards := [][]byte{
		{0x00, 0xFF, 0x7F},
		{0x80, 0x01, 0xFE},
		{0xAA, 0xBB, 0xCC},
		{0xDD, 0xEE, 0xFF},
	}

	err := store.Save("binary_account", originalShards)
	require.NoError(t, err)

	loadedShards, err := store.Load("binary_account", 4)

	require.NoError(t, err)
	assert.Equal(t, originalShards, loadedShards)
}

func TestFileStore_Integration_SplitSaveLoadReconstruct(t *testing.T) {
	tmpDir := t.TempDir()
	splitter, err := reedsolomon.NewSplitter(2, 2)
	require.NoError(t, err)
	store := NewFileStore(tmpDir)
	originalData := []byte("Integration test: split, save, load, reconstruct!")

	// Split
	shards, err := splitter.Split(originalData)
	require.NoError(t, err)

	// Save
	err = store.Save("integration_test", shards)
	require.NoError(t, err)

	// Simulate partial loss by deleting some files
	require.NoError(t, os.Remove(filepath.Join(tmpDir, "integration_test.1")))

	// Load (with missing shard)
	loadedShards, err := store.Load("integration_test", splitter.TotalShards())
	require.NoError(t, err)
	assert.Nil(t, loadedShards[1])

	// Reconstruct
	recovered, err := splitter.Reconstruct(loadedShards)
	require.NoError(t, err)
	assert.Equal(t, originalData, recovered)
}
