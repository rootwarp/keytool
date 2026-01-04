package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp_NameIsKeytool(t *testing.T) {
	app := NewApp()

	assert.Equal(t, "keytool", app.Name)
}

func TestNewApp_HasEthCommand(t *testing.T) {
	app := NewApp()

	var ethCmd *bool
	for _, cmd := range app.Commands {
		if cmd.Name == "eth" {
			found := true
			ethCmd = &found
			break
		}
	}

	assert.NotNil(t, ethCmd)
	assert.True(t, *ethCmd)
}
