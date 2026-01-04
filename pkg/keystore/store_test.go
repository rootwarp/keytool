package keystore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultStore_Save_CreatesKeystoreFile(t *testing.T) {
	store := NewDefaultStore()
	privateKey := "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
	tmpDir := t.TempDir()

	filePath, err := store.Save(privateKey, "testpassword", tmpDir)

	require.NoError(t, err)
	assert.FileExists(t, filePath)
}

func TestDefaultStore_Save_CreatesUniqueFilenames(t *testing.T) {
	store := NewDefaultStore()
	privateKey1 := "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
	privateKey2 := "0x5c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362319"
	tmpDir := t.TempDir()

	file1, err1 := store.Save(privateKey1, "testpassword1", tmpDir)
	file2, err2 := store.Save(privateKey2, "testpassword2", tmpDir)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, file1, file2)
}

func TestDefaultStore_Save_FileContainsEncryptedJSON(t *testing.T) {
	store := NewDefaultStore()
	privateKey := "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
	tmpDir := t.TempDir()

	filePath, err := store.Save(privateKey, "testpassword", tmpDir)

	require.NoError(t, err)

	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "crypto")
	assert.Contains(t, string(content), "address")
}

func TestDefaultStore_Save_CreatesDirectoryIfNotExists(t *testing.T) {
	store := NewDefaultStore()
	privateKey := "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
	tmpDir := filepath.Join(t.TempDir(), "nonexistent", "path")

	filePath, err := store.Save(privateKey, "testpassword", tmpDir)

	require.NoError(t, err)
	assert.FileExists(t, filePath)
}

func TestDefaultStore_Save_ErrorsOnInvalidPrivateKey(t *testing.T) {
	store := NewDefaultStore()
	tmpDir := t.TempDir()

	_, err := store.Save("invalid-key", "testpassword", tmpDir)

	assert.Error(t, err)
}

func TestDefaultStore_Save_ErrorsOnEmptyPassword(t *testing.T) {
	store := NewDefaultStore()
	privateKey := "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
	tmpDir := t.TempDir()

	_, err := store.Save(privateKey, "", tmpDir)

	assert.Error(t, err)
}
