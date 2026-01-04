package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNewAccountCmd_NameIsNewAccount(t *testing.T) {
	cmd := NewNewAccountCmd()

	assert.Equal(t, "new-account", cmd.Name)
}

func TestNewNewAccountCmd_HasKeystoreDirFlag(t *testing.T) {
	cmd := NewNewAccountCmd()

	var keystoreDirFlag *bool
	for _, flag := range cmd.Flags {
		if flag.Names()[0] == "keystore-dir" {
			found := true
			keystoreDirFlag = &found
			break
		}
	}

	assert.NotNil(t, keystoreDirFlag)
	assert.True(t, *keystoreDirFlag)
}

func TestNewNewAccountCmd_KeystoreDirFlagIsRequired(t *testing.T) {
	app := NewApp()

	err := app.Run([]string{"keytool", "eth", "new-account"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "keystore-dir")
}

func TestGenerateRandomIndex_ReturnsValueUnder1000000(t *testing.T) {
	for i := 0; i < 100; i++ {
		index, err := generateRandomIndex()

		assert.NoError(t, err)
		assert.Less(t, index, uint32(1000000))
	}
}

func TestZeroString_ClearsString(t *testing.T) {
	password := "sensitive-password"

	zeroString(&password)

	assert.Empty(t, password)
}

func TestZeroString_HandlesNil(t *testing.T) {
	assert.NotPanics(t, func() {
		zeroString(nil)
	})
}
