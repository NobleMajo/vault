<div align="center">

# 🛅 vault

![MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![CI/CD](https://github.com/NobleMajo/vault/actions/workflows/go-bin-release.yml/badge.svg)
![CI/CD](https://github.com/NobleMajo/vault/actions/workflows/go-test-build.yml/badge.svg)  
![](https://img.shields.io/badge/dynamic/json?color=green&label=watchers&query=watchers&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FNobleMajo%2Fvault)
![](https://img.shields.io/badge/dynamic/json?color=yellow&label=stars&query=stargazers_count&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FNobleMajo%2Fvault)
![](https://img.shields.io/badge/dynamic/json?color=navy&label=forks&query=forks&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FNobleMajo%2Fvault)

</div>

## About

Vault is a minimalistic CLI tool that encrypts and decrypts plain files into Vault files. (`.vt`).

The idea behind this tool is to have a CLI utility that can quickly and easily encrypt individual files, allowing users to securely store API tokens, secrets, credentials, or any private data on their own disk.

<details><summary><strong>Advertising</strong></summary>

### Advertising

_Are you also just a normal software developer or admin with lots of API keys, encryption keys or other secrets and credentials?_
_Or do you simply have logs or plain text files that you want to send to someone securely?_
**Then I have exactly what you are looking for today!**

_Hold on tight and take a closer look at this cli tool, because it might meet your exact needs._

</details>

<details><summary><strong>Encryption</strong></summary>

### Encryption

Vault uses asymmetric RSA encryption and symmetric AES-256 encryption to keep your data as secure as possible.
To do this, vault uses private and public key on disk (default: `~/.ssh/id_rsa.pub`) and also asks you for a password.

Currently no elliptic curve support! Just rsa.

</details>

<details><summary><strong>Usage</strong></summary>

# Usage

Vault operations are sub commands defined via the first command line argument.

## Help

The following block is the main help output if you do not use a subcommand or use help:

```ts
Vault is a file encryption and decryption cli tool written in go.
For more help, visit https://github.com/NobleMajo/vault

Usage:
  vault [flags]
  vault [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Create a initial encrypted vault file for default text
  lock        Locks your plain file into a vault file
  passwd      Changes the password of your vault file
  print       Prints the decrypted content of your vault file
  temp        Temporary unlocks your vault file into a plain file
  unlock      Unlocks your vault file into a plain file
  version     Prints version message

Flags:
  -h, --help      help for vault
  -b, --verbose   enable verbose mode (VAULT_VERBOSE)
  -v, --version   prints version

Use "vault [command] --help" for more information about a command.
```

### init

Create a new locked vault file:

```sh
vault init
```

**OR**

### lock

Add some content to your `vault.txt` and lock it:

```sh
vim vault.txt
vault lock
```

### unlock

Unlock the vault as plain `.txt` file:

```sh
vault unlock
```

### temp

Unlock the file for 5 seconds as `.txt`.
In this time you can open it with an editor.

```sh
vault temp
```

### print

Print the locked content in console:

```sh
vault print
```

## Other filename

To choose another file than the `vault.txt` use the second argument without extensions:
(`test` for `test.txt` and `test.vt`)

```sh
vault lock <filename>
vault temp <filename>
vault unlock <filename>
vault init <filename>
vault print <filename>
```

</details>

<details><summary><strong>User Guide</strong></summary>

# User Guide

## Requirements

Linux- or macos-like systems with `go` or `wget & tar` installed.

## Getting Started

Start the latest repo version directly without leaving stuff in the current working dir:

```sh
go run github.com/NobleMajo/vault@latest
```

## Quick help

```sh
go run github.com/NobleMajo/vault@latest -h
```

## Install via go

###### _For this section go is required, check out the [install go guide](#install-go)._

```sh
go install github.com/NobleMajo/vault@latest
```

## Install via wget

```sh
export CUSTOM_BIN_DIR="/usr/local/bin" # <- change if needed
export CUSTOM_VERSION="" # <- set latest version here

rm -rf $CUSTOM_BIN_DIR/vault
wget https://github.com/NobleMajo/vault/releases/download/v$CUSTOM_VERSION/vault-v$CUSTOM_VERSION-linux-amd64.tar.gz -O /tmp/vault.tar.gz
tar -xzvf /tmp/vault.tar.gz -C $CUSTOM_BIN_DIR/ vault
rm /tmp/vault.tar.gz
```

# Build

## Build requirements

To build, you need to install go.
The required go version is in the `go.mod` file.

## Build Instructions

###### _For this section go is required, check out the [install go guide](#install-go)._

Clone the repo:

```sh
git clone https://github.com/NobleMajo/vault.git
cd vault
```

Build the vault binary from source code:

```sh
make build
./vault
```

</details>

<details><summary><strong>Development</strong></summary>

# Development

###### _For this section go is required, check out the [install go guide](#install-go)._

This part is work in progress, I want to use 'AIR' as auto-reload tool:

```sh
make dev #WIP
```

## Install go

The required go version for this project is in the `go.mod` file.

To install and update go, I can recommend the following repo:

```sh
git clone git@github.com:udhos/update-golang.git golang-updater
cd golang-updater
sudo ./update-golang.sh
```

</details>

<div align="center">

# 🤝 Contributing

Contributions to this project are welcome!  
Follow the [CONTRIBUTING.md](CONTRIBUTING.md) for more infos.

# ⚠️ Disclaimer

This project is provided without warranties.

# 📜 License

Licensed under the [MIT license](LICENSE).

<a href="https://discord.coreunit.net">
    <img alt="CoreUnit.NET Discord Banner" src="https://discord.com/api/guilds/422136748294930443/widget.png?style=banner2">
</a>

</div>