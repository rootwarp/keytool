// Package secure provides utilities for clearing sensitive data from memory.
package secure

// ZeroBytes overwrites a byte slice with zeros to clear sensitive data from memory.
// This helps prevent sensitive information like private keys from lingering in memory.
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// ZeroString clears a string by setting it to empty and returns an empty string.
// Since Go strings are immutable, this replaces the pointer's value with an empty string.
// The original string's memory may still contain data until garbage collected,
// but the reference is cleared immediately.
func ZeroString(s *string) string {
	if s == nil {
		return ""
	}
	*s = ""
	return ""
}
