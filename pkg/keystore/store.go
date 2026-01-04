// Package keystore provides Ethereum keystore persistence functionality.
package keystore

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

// Store defines the interface for account persistence.
type Store interface {
	// Save persists a private key to the keystore with the given password.
	// Returns the path to the created keystore file.
	Save(privateKeyHex string, password string, keystorePath string) (string, error)
}

// DefaultStore implements Store using go-ethereum keystore.
type DefaultStore struct{}

// NewDefaultStore creates a new DefaultStore instance.
func NewDefaultStore() *DefaultStore {
	return &DefaultStore{}
}

// Save persists a private key to the keystore with the given password.
func (s *DefaultStore) Save(privateKeyHex, password, keystorePath string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	if err := os.MkdirAll(keystorePath, 0700); err != nil {
		return "", fmt.Errorf("failed to create keystore directory: %w", err)
	}

	ks := keystore.NewKeyStore(keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := ks.ImportECDSA(privateKey, password)
	if err != nil {
		return "", fmt.Errorf("failed to import key to keystore: %w", err)
	}

	return account.URL.Path, nil
}
