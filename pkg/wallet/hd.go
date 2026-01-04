// Package wallet provides HD wallet derivation for Ethereum accounts.
package wallet

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
)

// Account represents a derived Ethereum account.
type Account struct {
	Address    string `json:"address"`
	HDPath     string `json:"hd_path"`
	Mnemonic   string `json:"mnemonic"`
	PrivateKey string `json:"private_key"`
}

// Zero clears sensitive fields (Mnemonic and PrivateKey) from memory.
// Should be called when the account data is no longer needed.
func (a *Account) Zero() {
	a.Mnemonic = ""
	a.PrivateKey = ""
}

// HDDeriver defines the interface for HD wallet derivation.
type HDDeriver interface {
	// DeriveAccount derives an Ethereum account from mnemonic and HD path.
	DeriveAccount(mnemonic string, hdPath string) (*Account, error)

	// GenerateHDPath generates a valid BIP-44 HD path with given index.
	// Path format: m/44'/60'/0'/0/{index}
	GenerateHDPath(index uint32) string
}

// DefaultHDDeriver implements HDDeriver using go-ethereum-hdwallet.
type DefaultHDDeriver struct{}

// NewDefaultHDDeriver creates a new DefaultHDDeriver instance.
func NewDefaultHDDeriver() *DefaultHDDeriver {
	return &DefaultHDDeriver{}
}

// GenerateHDPath generates a valid BIP-44 HD path for Ethereum with given index.
func (d *DefaultHDDeriver) GenerateHDPath(index uint32) string {
	return fmt.Sprintf("m/44'/60'/0'/0/%d", index)
}

// DeriveAccount derives an Ethereum account from mnemonic and HD path.
func (d *DefaultHDDeriver) DeriveAccount(mnemonic string, hdPath string) (*Account, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("invalid mnemonic phrase: %w", err)
	}

	path, err := hdwallet.ParseDerivationPath(hdPath)
	if err != nil {
		return nil, fmt.Errorf("invalid HD path: %w", err)
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %w", err)
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	return &Account{
		Address:    account.Address.Hex(),
		HDPath:     hdPath,
		Mnemonic:   mnemonic,
		PrivateKey: "0x" + hex.EncodeToString(privateKeyBytes),
	}, nil
}
