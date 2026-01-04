package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEthCmd_NameIsEth(t *testing.T) {
	cmd := NewEthCmd()

	assert.Equal(t, "eth", cmd.Name)
}

func TestNewEthCmd_HasNewAccountSubcommand(t *testing.T) {
	cmd := NewEthCmd()

	var newAccountCmd *bool
	for _, sub := range cmd.Subcommands {
		if sub.Name == "new-account" {
			found := true
			newAccountCmd = &found
			break
		}
	}

	assert.NotNil(t, newAccountCmd)
	assert.True(t, *newAccountCmd)
}
