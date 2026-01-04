package cli

import (
	"github.com/urfave/cli/v2"
)

// NewEthCmd creates the eth subcommand.
func NewEthCmd() *cli.Command {
	return &cli.Command{
		Name:  "eth",
		Usage: "Ethereum-related commands",
		Description: `Commands for managing Ethereum keys and accounts.`,
		Subcommands: []*cli.Command{
			NewNewAccountCmd(),
		},
	}
}
