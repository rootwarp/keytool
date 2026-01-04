# keytool

Ethereum key management tool for creating and managing accounts with BIP-39/BIP-44 HD wallet support.

## Build

```bash
make build
```

Binary will be created at `bin/keytool`.

## Usage

### Create a New Account

```bash
keytool eth new-account --keystore-dir <path>
```

Creates a new Ethereum account with:
- BIP-39 mnemonic (24 words)
- BIP-44 HD derivation path (`m/44'/60'/0'/0/{index}`)
- Encrypted keystore file

#### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--keystore-dir` | Yes | Directory to save the encrypted keystore file |
| `--shard-dir` | No | Directory to save secret shards (enables Reed-Solomon splitting) |

#### Examples

Basic usage:
```bash
keytool eth new-account --keystore-dir ./keystore
```

With secret splitting:
```bash
keytool eth new-account --keystore-dir ./keystore --shard-dir ./shards
```
