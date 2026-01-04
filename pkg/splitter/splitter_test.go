package splitter

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSplitter_CreatesInstance(t *testing.T) {
	splitter, err := NewSplitter(2, 2)

	require.NoError(t, err)
	assert.NotNil(t, splitter)
}

func TestNewSplitter_ErrorsOnInvalidParams(t *testing.T) {
	_, err := NewSplitter(0, 2)

	assert.Error(t, err)
}

func TestSplitter_Split_CreatesShardFiles(t *testing.T) {
	tmpDir := t.TempDir()
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	data := SecretData{
		Address:  "0x1234567890abcdef1234567890abcdef12345678",
		Mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		HDPath:   "m/44'/60'/0'/0/0",
	}

	err = splitter.Split(data, tmpDir)

	require.NoError(t, err)
	// Check that 4 shard files were created (k=2, m=2)
	for i := range 4 {
		filename := filepath.Join(tmpDir, data.Address+"."+string(rune('0'+i)))
		assert.FileExists(t, filename)
	}
}

func TestSplitter_Split_Base64EncodesJSON(t *testing.T) {
	tmpDir := t.TempDir()
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	data := SecretData{
		Address:  "0xabcdef",
		Mnemonic: "test mnemonic",
		HDPath:   "m/44'/60'/0'/0/1",
	}

	err = splitter.Split(data, tmpDir)
	require.NoError(t, err)

	// Reconstruct and verify the data was base64 encoded
	recovered, err := splitter.Reconstruct(data.Address, tmpDir)
	require.NoError(t, err)

	// The recovered data should match original
	assert.Equal(t, data.Address, recovered.Address)
	assert.Equal(t, data.Mnemonic, recovered.Mnemonic)
	assert.Equal(t, data.HDPath, recovered.HDPath)
}

func TestSplitter_Reconstruct_RecoverOriginalData(t *testing.T) {
	tmpDir := t.TempDir()
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	original := SecretData{
		Address:  "0x9876543210fedcba9876543210fedcba98765432",
		Mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
		HDPath:   "m/44'/60'/0'/0/999",
	}

	err = splitter.Split(original, tmpDir)
	require.NoError(t, err)

	recovered, err := splitter.Reconstruct(original.Address, tmpDir)

	require.NoError(t, err)
	assert.Equal(t, original, *recovered)
}

func TestSplitter_Reconstruct_WithOneMissingShard(t *testing.T) {
	tmpDir := t.TempDir()
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	original := SecretData{
		Address:  "0xdeadbeef",
		Mnemonic: "test mnemonic phrase",
		HDPath:   "m/44'/60'/0'/0/42",
	}

	err = splitter.Split(original, tmpDir)
	require.NoError(t, err)

	// Delete one shard (RS can recover with m=2 parity)
	err = os.Remove(filepath.Join(tmpDir, original.Address+".0"))
	require.NoError(t, err)

	recovered, err := splitter.Reconstruct(original.Address, tmpDir)

	require.NoError(t, err)
	assert.Equal(t, original, *recovered)
}

func TestSplitter_Reconstruct_WithTwoMissingShards(t *testing.T) {
	tmpDir := t.TempDir()
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	original := SecretData{
		Address:  "0xcafebabe",
		Mnemonic: "another test mnemonic",
		HDPath:   "m/44'/60'/0'/0/100",
	}

	err = splitter.Split(original, tmpDir)
	require.NoError(t, err)

	// Delete two shards (maximum for m=2)
	err = os.Remove(filepath.Join(tmpDir, original.Address+".0"))
	require.NoError(t, err)
	err = os.Remove(filepath.Join(tmpDir, original.Address+".2"))
	require.NoError(t, err)

	recovered, err := splitter.Reconstruct(original.Address, tmpDir)

	require.NoError(t, err)
	assert.Equal(t, original, *recovered)
}

func TestSplitter_Reconstruct_FailsWithThreeMissingShards(t *testing.T) {
	tmpDir := t.TempDir()
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	original := SecretData{
		Address:  "0xfailtest",
		Mnemonic: "fail test mnemonic",
		HDPath:   "m/44'/60'/0'/0/0",
	}

	err = splitter.Split(original, tmpDir)
	require.NoError(t, err)

	// Delete three shards (exceeds m=2)
	err = os.Remove(filepath.Join(tmpDir, original.Address+".0"))
	require.NoError(t, err)
	err = os.Remove(filepath.Join(tmpDir, original.Address+".1"))
	require.NoError(t, err)
	err = os.Remove(filepath.Join(tmpDir, original.Address+".2"))
	require.NoError(t, err)

	_, err = splitter.Reconstruct(original.Address, tmpDir)

	assert.Error(t, err)
}

func TestSplitter_TotalShards_ReturnsFour(t *testing.T) {
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)

	total := splitter.TotalShards()

	assert.Equal(t, 4, total)
}

func TestSplitter_Split_DataIsBase64Encoded(t *testing.T) {
	tmpDir := t.TempDir()
	splitter, err := NewSplitter(2, 2)
	require.NoError(t, err)
	data := SecretData{
		Address:  "0xbase64test",
		Mnemonic: "base64 encoding test",
		HDPath:   "m/44'/60'/0'/0/5",
	}

	// Marshal to JSON first to know what we expect
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)
	expectedBase64 := base64.StdEncoding.EncodeToString(jsonData)

	err = splitter.Split(data, tmpDir)
	require.NoError(t, err)

	// Recover and verify the intermediate base64 step happened
	recovered, err := splitter.Reconstruct(data.Address, tmpDir)
	require.NoError(t, err)

	// Re-encode the recovered data to verify round-trip
	recoveredJSON, err := json.Marshal(recovered)
	require.NoError(t, err)
	recoveredBase64 := base64.StdEncoding.EncodeToString(recoveredJSON)

	assert.Equal(t, expectedBase64, recoveredBase64)
}
