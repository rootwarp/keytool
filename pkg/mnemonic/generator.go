// Package mnemonic provides BIP-39 mnemonic generation and validation.
package mnemonic

import (
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

// Generator defines the interface for mnemonic generation.
type Generator interface {
	// Generate creates a new BIP-39 mnemonic with the specified bit size.
	// bitSize should be 128, 160, 192, 224, or 256 (128 = 12 words, 256 = 24 words).
	Generate(bitSize int) (string, error)

	// Validate checks if a mnemonic is valid according to BIP-39.
	Validate(mnemonic string) bool
}

// DefaultGenerator implements Generator using go-bip39.
type DefaultGenerator struct{}

// NewDefaultGenerator creates a new DefaultGenerator instance.
func NewDefaultGenerator() *DefaultGenerator {
	return &DefaultGenerator{}
}

// Generate creates a new BIP-39 mnemonic with the specified bit size.
func (g *DefaultGenerator) Generate(bitSize int) (string, error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// Validate checks if a mnemonic is valid according to BIP-39.
func (g *DefaultGenerator) Validate(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}
