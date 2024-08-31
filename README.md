# Vault
![CI/CD](https://github.com/noblemajo/vault/actions/workflows/go-bin-release.yml/badge.svg)
![CI/CD](https://github.com/noblemajo/vault/actions/workflows/go-test-build.yml/badge.svg)  
![MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/dynamic/json?color=green&label=watchers&query=watchers&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fvault)
![](https://img.shields.io/badge/dynamic/json?color=yellow&label=stars&query=stargazers_count&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fvault)
![](https://img.shields.io/badge/dynamic/json?color=navy&label=forks&query=forks&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fvault)

Vault is a minimalistic CLI tool that encrypts and decrypts plain files into Vault files. (`.vt`).

The idea behind this tool is to have a CLI utility that can quickly and easily encrypt individual files, allowing users to securely store API tokens, secrets, credentials, or any private data on their own disk.

# Table of Contents
- [Vault](#vault)
- [Table of Contents](#table-of-contents)
  - [Advertising](#advertising)
  - [Encryption](#encryption)
- [Getting Started](#getting-started)
  - [Requirements](#requirements)
  - [Install via go](#install-via-go)
  - [Install via wget](#install-via-wget)
- [Build](#build)
  - [Build requirements](#build-requirements)
- [Usage](#usage)
  - [Help](#help)
    - [init](#init)
    - [lock](#lock)
    - [unlock](#unlock)
    - [temp](#temp)
    - [print](#print)
  - [Other filename](#other-filename)
  - [Build](#build-1)
- [Development](#development)
  - [Install go](#install-go)
- [Contributing](#contributing)
- [License](#license)
- [Disclaimer](#disclaimer)

## Advertising
*Are you also just a normal software developer or admin with lots of API keys, encryption keys or other secrets and credentials?*
*Or do you simply have logs or plain text files that you want to send to someone securely?*
**Then I have exactly what you are looking for today!**

*Hold on tight and take a closer look at this command line interface tool because it might meet your exact needs.*

## Encryption
Vault uses asymmetric RSA encryption and symmetric AES-256 encryption to keep your data as secure as possible.
To do this, vault uses private and public key on disk (default: `~/.ssh/id_rsa.pub`) and also asks you for a password.

Currently no elliptic curve support! Just rsa.

# Getting Started

## Requirements
None windows system with `go` or `wget & tar` installed.

## Install via go
###### *For this section go is required, check out the [install go guide](#install-go).*

```sh
go install https://github.com/NobleMajo/vault
```

## Install via wget
```sh
BIN_DIR="/usr/local/bin"
VAULT_VERSION="1.3.3"

rm -rf $BIN_DIR/vault
wget https://github.com/NobleMajo/vault/releases/download/v$VAULT_VERSION/vault-v$VAULT_VERSION-linux-amd64.tar.gz -O /tmp/vault.tar.gz
tar -xzvf /tmp/vault.tar.gz -C $BIN_DIR/ vault
rm /tmp/vault.tar.gz
```

# Build
## Build requirements
To build, you need to install go. 
The required go version is in the `go.mod` file.

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
To choose a other file then the `vault.txt` use the second argument without extensions:
(`test` for `test.txt` and `test.vt`)
 ```sh
vault lock <filename>
vault temp <filename>
vault unlock <filename>
vault init <filename>
vault print <filename>
```

## Build
###### *For this section go is required, check out the [install go guide](#install-go).*

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

# Development
###### *For this section go is required, check out the [install go guide](#install-go).*

This part is work in process, i want use 'AIR' as autoreload tool:
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

# Contributing
Contributions to Vault are welcome!  
Interested users can refer to the guidelines provided in the [CONTRIBUTING.md](CONTRIBUTING.md) file to contribute to the project and help improve its functionality and features.

# License
Vault is licensed under the [MIT license](LICENSE), providing users with flexibility and freedom to use and modify the software according to their needs.

# Disclaimer
Vault is provided without warranties.  
Users are advised to review the accompanying license for more information on the terms of use and limitations of liability.
