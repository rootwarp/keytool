// Package splitter provides functionality for splitting secrets using Reed-Solomon encoding.
package splitter

// SecretData contains the account information to be split and distributed.
type SecretData struct {
	Address  string `json:"address"`
	Mnemonic string `json:"mnemonic"`
	HDPath   string `json:"hd_path"`
}
