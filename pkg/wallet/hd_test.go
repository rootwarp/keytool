package wallet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultHDDeriver_GenerateHDPath_ReturnsValidBIP44Path(t *testing.T) {
	deriver := NewDefaultHDDeriver()

	tests := []struct {
		index    uint32
		expected string
	}{
		{0, "m/44'/60'/0'/0/0"},
		{1, "m/44'/60'/0'/0/1"},
		{5, "m/44'/60'/0'/0/5"},
		{999999, "m/44'/60'/0'/0/999999"},
	}

	for _, tt := range tests {
		path := deriver.GenerateHDPath(tt.index)
		assert.Equal(t, tt.expected, path)
	}
}

func TestDefaultHDDeriver_DeriveAccount_ReturnsValidAccount(t *testing.T) {
	deriver := NewDefaultHDDeriver()
	testMnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon " +
		"abandon abandon abandon about"

	account, err := deriver.DeriveAccount(testMnemonic, "m/44'/60'/0'/0/0")

	require.NoError(t, err)
	assert.NotEmpty(t, account.Address)
	assert.Equal(t, "m/44'/60'/0'/0/0", account.HDPath)
	assert.Equal(t, testMnemonic, account.Mnemonic)
	assert.NotEmpty(t, account.PrivateKey)
	assert.True(t, strings.HasPrefix(account.Address, "0x"))
	assert.True(t, strings.HasPrefix(account.PrivateKey, "0x"))
}

func TestDefaultHDDeriver_DeriveAccount_DeterministicResults(t *testing.T) {
	deriver := NewDefaultHDDeriver()
	testMnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon " +
		"abandon abandon abandon about"

	account1, err1 := deriver.DeriveAccount(testMnemonic, "m/44'/60'/0'/0/0")
	account2, err2 := deriver.DeriveAccount(testMnemonic, "m/44'/60'/0'/0/0")

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, account1.Address, account2.Address)
	assert.Equal(t, account1.PrivateKey, account2.PrivateKey)
}

func TestDefaultHDDeriver_DeriveAccount_DifferentIndexesDifferentAddresses(t *testing.T) {
	deriver := NewDefaultHDDeriver()
	testMnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon " +
		"abandon abandon abandon about"

	account1, err1 := deriver.DeriveAccount(testMnemonic, "m/44'/60'/0'/0/0")
	account2, err2 := deriver.DeriveAccount(testMnemonic, "m/44'/60'/0'/0/1")

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, account1.Address, account2.Address)
	assert.NotEqual(t, account1.PrivateKey, account2.PrivateKey)
}

func TestDefaultHDDeriver_DeriveAccount_ErrorsOnInvalidMnemonic(t *testing.T) {
	deriver := NewDefaultHDDeriver()

	_, err := deriver.DeriveAccount("invalid mnemonic words", "m/44'/60'/0'/0/0")

	assert.Error(t, err)
}

func TestDefaultHDDeriver_DeriveAccount_ErrorsOnInvalidHDPath(t *testing.T) {
	deriver := NewDefaultHDDeriver()
	testMnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon " +
		"abandon abandon abandon about"

	_, err := deriver.DeriveAccount(testMnemonic, "invalid/path")

	assert.Error(t, err)
}

func TestAccount_Zero_ClearsSensitiveFields(t *testing.T) {
	account := &Account{
		Address:    "0x1234567890abcdef",
		HDPath:     "m/44'/60'/0'/0/0",
		Mnemonic:   "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		PrivateKey: "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318",
	}

	account.Zero()

	assert.Empty(t, account.Mnemonic, "Mnemonic should be cleared")
	assert.Empty(t, account.PrivateKey, "PrivateKey should be cleared")
	assert.Equal(t, "0x1234567890abcdef", account.Address, "Address should remain")
	assert.Equal(t, "m/44'/60'/0'/0/0", account.HDPath, "HDPath should remain")
}

func TestAccount_Zero_HandlesEmptyFields(t *testing.T) {
	account := &Account{}

	assert.NotPanics(t, func() {
		account.Zero()
	})
}
