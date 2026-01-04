// Package cli provides command-line interface commands for keytool.
package cli

import (
	"github.com/urfave/cli/v2"
)

// NewApp creates the CLI application for keytool.
func NewApp() *cli.App {
	return &cli.App{
		Name:  "keytool",
		Usage: "A tool for managing cryptocurrency keys",
		Description: `keytool is a CLI tool for creating and managing
cryptocurrency keys and accounts.`,
		Commands: []*cli.Command{
			NewEthCmd(),
		},
	}
}
