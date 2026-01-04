package reedsolomon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSplitter_CreatesWithValidParams(t *testing.T) {
	splitter, err := NewSplitter(2, 2)

	require.NoError(t, err)
	assert.NotNil(t, splitter)
}

func TestNewSplitter_ErrorsOnInvalidParams(t *testing.T) {
	_, err := NewSplitter(0, 2)

	assert.Error(t, err)
}

func TestSplitter_TotalShards_ReturnsDataPlusParity(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)

	total := splitter.TotalShards()

	assert.Equal(t, 4, total)
}

func TestSplitter_Split_ReturnsFourShards(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	data := []byte("test data for splitting")

	shards, err := splitter.Split(data)

	require.NoError(t, err)
	assert.Len(t, shards, 4)
}

func TestSplitter_Split_AllShardsHaveData(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	data := []byte("test data for splitting")

	shards, err := splitter.Split(data)

	require.NoError(t, err)
	for i, shard := range shards {
		assert.NotEmpty(t, shard, "shard %d should not be empty", i)
	}
}

func TestSplitter_Reconstruct_RecoverOriginalData(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	originalData := []byte("Hello, secret splitting test!")

	shards, err := splitter.Split(originalData)
	require.NoError(t, err)

	recovered, err := splitter.Reconstruct(shards)

	require.NoError(t, err)
	assert.Equal(t, originalData, recovered)
}

func TestSplitter_Reconstruct_RecoverWithOneMissingShard(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	originalData := []byte("Recover with missing shard")

	shards, err := splitter.Split(originalData)
	require.NoError(t, err)

	// Simulate loss of one shard
	shards[0] = nil

	recovered, err := splitter.Reconstruct(shards)

	require.NoError(t, err)
	assert.Equal(t, originalData, recovered)
}

func TestSplitter_Reconstruct_RecoverWithTwoMissingShards(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	originalData := []byte("Recover with two missing shards")

	shards, err := splitter.Split(originalData)
	require.NoError(t, err)

	// Simulate loss of two shards (maximum for m=2)
	shards[0] = nil
	shards[2] = nil

	recovered, err := splitter.Reconstruct(shards)

	require.NoError(t, err)
	assert.Equal(t, originalData, recovered)
}

func TestSplitter_Reconstruct_FailsWithThreeMissingShards(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	originalData := []byte("Cannot recover with three missing")

	shards, err := splitter.Split(originalData)
	require.NoError(t, err)

	// Simulate loss of three shards (exceeds m=2)
	shards[0] = nil
	shards[1] = nil
	shards[2] = nil

	_, err = splitter.Reconstruct(shards)

	assert.Error(t, err)
}

func TestSplitter_SplitReconstruct_HandlesSmallData(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	originalData := []byte("X") // Single byte

	shards, err := splitter.Split(originalData)
	require.NoError(t, err)

	recovered, err := splitter.Reconstruct(shards)

	require.NoError(t, err)
	assert.Equal(t, originalData, recovered)
}

func TestSplitter_SplitReconstruct_HandlesBinaryData(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	originalData := []byte{0x00, 0xFF, 0x7F, 0x80, 0x01, 0xFE}

	shards, err := splitter.Split(originalData)
	require.NoError(t, err)

	recovered, err := splitter.Reconstruct(shards)

	require.NoError(t, err)
	assert.Equal(t, originalData, recovered)
}
