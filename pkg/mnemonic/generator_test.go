package mnemonic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenerator_Generate_Returns12WordMnemonic(t *testing.T) {
	gen := NewDefaultGenerator()
	mnemonic, err := gen.Generate(128)

	require.NoError(t, err)
	words := strings.Split(mnemonic, " ")
	assert.Len(t, words, 12)
}

func TestDefaultGenerator_Generate_Returns24WordMnemonic(t *testing.T) {
	gen := NewDefaultGenerator()
	mnemonic, err := gen.Generate(256)

	require.NoError(t, err)
	words := strings.Split(mnemonic, " ")
	assert.Len(t, words, 24)
}

func TestDefaultGenerator_Generate_ReturnsDifferentMnemonicsEachCall(t *testing.T) {
	gen := NewDefaultGenerator()
	m1, err1 := gen.Generate(128)
	m2, err2 := gen.Generate(128)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, m1, m2)
}

func TestDefaultGenerator_Generate_ErrorsOnInvalidBitSize(t *testing.T) {
	gen := NewDefaultGenerator()
	_, err := gen.Generate(100)

	assert.Error(t, err)
}

func TestDefaultGenerator_Validate_ReturnsTrueForValidMnemonic(t *testing.T) {
	gen := NewDefaultGenerator()
	mnemonic, err := gen.Generate(128)

	require.NoError(t, err)
	assert.True(t, gen.Validate(mnemonic))
}

func TestDefaultGenerator_Validate_ReturnsFalseForInvalidMnemonic(t *testing.T) {
	gen := NewDefaultGenerator()

	assert.False(t, gen.Validate("invalid mnemonic words that are not real"))
}

func TestDefaultGenerator_Validate_ReturnsTrueForKnownTestMnemonic(t *testing.T) {
	gen := NewDefaultGenerator()
	testMnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon " +
		"abandon abandon abandon about"

	assert.True(t, gen.Validate(testMnemonic))
}
