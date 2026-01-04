package secure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroBytes_ClearsAllBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	ZeroBytes(data)

	for i, b := range data {
		assert.Equal(t, byte(0), b, "byte at index %d should be zero", i)
	}
}

func TestZeroBytes_HandlesEmptySlice(t *testing.T) {
	data := []byte{}

	ZeroBytes(data)

	assert.Empty(t, data)
}

func TestZeroBytes_HandlesNilSlice(t *testing.T) {
	var data []byte

	assert.NotPanics(t, func() {
		ZeroBytes(data)
	})
}

func TestZeroString_ReturnsEmptyString(t *testing.T) {
	original := "sensitive mnemonic phrase here"

	result := ZeroString(&original)

	assert.Empty(t, result)
	assert.Empty(t, original)
}

func TestZeroString_HandlesEmptyString(t *testing.T) {
	original := ""

	result := ZeroString(&original)

	assert.Empty(t, result)
}

func TestZeroString_HandlesNilPointer(t *testing.T) {
	assert.NotPanics(t, func() {
		result := ZeroString(nil)
		assert.Empty(t, result)
	})
}
