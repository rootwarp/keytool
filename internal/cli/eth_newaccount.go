package cli

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rootwarp/keytool/pkg/keystore"
	"github.com/rootwarp/keytool/pkg/mnemonic"
	"github.com/rootwarp/keytool/pkg/splitter"
	"github.com/rootwarp/keytool/pkg/wallet"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

// NewNewAccountCmd creates the new-account subcommand.
func NewNewAccountCmd() *cli.Command {
	return &cli.Command{
		Name:  "new-account",
		Usage: "Create a new Ethereum account",
		Description: `Create a new Ethereum account with a BIP-39 mnemonic and BIP-44 derivation path.

The account information will be printed as JSON, and the keystore file will be
saved to the specified directory.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "keystore-dir",
				Usage:    "Directory to save keystore file (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "shard-dir",
				Usage: "Directory to save secret shards (optional, enables secret splitting)",
			},
		},
		Action: runNewAccount,
	}
}

func runNewAccount(c *cli.Context) error {
	keystoreDir := c.String("keystore-dir")

	mnemonicGen := mnemonic.NewDefaultGenerator()
	hdDeriver := wallet.NewDefaultHDDeriver()
	ks := keystore.NewDefaultStore()

	mnemonicPhrase, err := mnemonicGen.Generate(256)
	if err != nil {
		return fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	index, err := generateRandomIndex()
	if err != nil {
		return fmt.Errorf("failed to generate random index: %w", err)
	}

	hdPath := hdDeriver.GenerateHDPath(index)

	account, err := hdDeriver.DeriveAccount(mnemonicPhrase, hdPath)
	if err != nil {
		return fmt.Errorf("failed to derive account: %w", err)
	}
	defer account.Zero()

	output, err := json.MarshalIndent(account, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal account: %w", err)
	}

	_, _ = fmt.Fprintln(c.App.Writer, string(output))

	password, err := promptPassword(c)
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	defer zeroString(&password)

	keystorePath, err := ks.Save(account.PrivateKey, password, keystoreDir)
	if err != nil {
		return fmt.Errorf("failed to save keystore: %w", err)
	}

	_, _ = fmt.Fprintf(c.App.Writer, "\nKeystore saved to: %s\n", keystorePath)

	// Split secret if shard-dir is provided
	shardDir := c.String("shard-dir")
	if shardDir != "" {
		secretSplitter, err := splitter.NewSplitter(2, 2)
		if err != nil {
			return fmt.Errorf("failed to create splitter: %w", err)
		}

		secretData := splitter.SecretData{
			Address:  account.Address,
			Mnemonic: account.Mnemonic,
			HDPath:   account.HDPath,
		}

		if err := secretSplitter.Split(secretData, shardDir); err != nil {
			return fmt.Errorf("failed to split secret: %w", err)
		}

		_, _ = fmt.Fprintf(c.App.Writer, "Secret shards saved to: %s (%d files)\n",
			shardDir, secretSplitter.TotalShards())
	}

	return nil
}

// zeroString clears a string by setting it to empty.
func zeroString(s *string) {
	if s != nil {
		*s = ""
	}
}

func generateRandomIndex() (uint32, error) {
	var buf [4]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return 0, err
	}
	index := binary.BigEndian.Uint32(buf[:]) % 1000000
	return index, nil
}

func promptPassword(c *cli.Context) (string, error) {
	_, err := fmt.Fprint(c.App.ErrWriter, "\nEnter password for keystore: ")
	if err != nil {
		return "", err
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		_, _ = fmt.Fprintln(c.App.ErrWriter)
		if err != nil {
			return "", err
		}

		password := strings.TrimSpace(string(passwordBytes))
		if password == "" {
			return "", fmt.Errorf("password cannot be empty")
		}

		_, _ = fmt.Fprint(c.App.ErrWriter, "Confirm password: ")
		confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		_, _ = fmt.Fprintln(c.App.ErrWriter)
		if err != nil {
			return "", err
		}

		if string(confirmBytes) != password {
			return "", fmt.Errorf("passwords do not match")
		}

		return password, nil
	}

	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	password = strings.TrimSpace(password)
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	return password, nil
}
