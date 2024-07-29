# Vault
![CI/CD](https://github.com/noblemajo/vault/actions/workflows/go-bin-release.yml/badge.svg)
![CI/CD](https://github.com/noblemajo/vault/actions/workflows/go-test-build.yml/badge.svg)  
![MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/dynamic/json?color=green&label=watchers&query=watchers&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fvault)
![](https://img.shields.io/badge/dynamic/json?color=yellow&label=stars&query=stargazers_count&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fvault)
![](https://img.shields.io/badge/dynamic/json?color=navy&label=forks&query=forks&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fvault)

Vault is a small and simple CLI tool that encrypt and decrypt plain files into vault-files (`.vt`).

The idea behind this tool is to have a CLI utility that can quickly and easily encrypt individual files, allowing users to securely store API tokens, secrets, credentials, or any private data on their own disk.

# Table of Contents
- [Vault](#vault)
- [Table of Contents](#table-of-contents)
  - [Advertising](#advertising)
  - [Encryption](#encryption)
- [Getting Started](#getting-started)
  - [Requirements](#requirements)
  - [Build](#build)
  - [Help](#help)
- [Usage](#usage)
  - [Other filename](#other-filename)
- [Development](#development)
  - [Automatic building](#automatic-building)
  - [Global binary linking](#global-binary-linking)
- [Contributing](#contributing)
- [License](#license)
- [Disclaimer](#disclaimer)

## Advertising
*Are you also a small software developer or admin with lots of API keys, encryption keys or other secrets and credentials?*
*Or do you simply have logs or plain text files that you want to send to someone securely?*
**Then I have exactly what you are looking for today!**

*Hold on tight and take a closer look at this command line interface tool because it might meet your exact needs.*

## Encryption
Vault uses asymmetric RSA encryption and symmetric AES-256 encryption to keep your data as secure as possible.
To do this, vault uses private and public key on disk (default: `~/.ssh/id_rsa.pub`) and also asks you for a password.

Currently no elliptic curve support! Just rsa.

# Getting Started
## Requirements
For building you need to install go.
For that i can recommend the following repo:
```sh
git clone git@github.com:udhos/update-golang.git golang-updater
cd golang-updater
sudo ./update-golang.sh
```

## Build
Clone the repo:
```sh
git clone https://github.com/NobleMajo/vault.git
cd vault
```

Build the vault binary from source code:
```sh
go build -o bin/vault cmd/main.go
# or npm run build
```

## Help
Run the help command on the binary
```sh
./vault
# or npm i -g .
# and then use "vault" anywhere on the machine
```

Output:
```rust
Usage:  ./vault [OPTIONS] COMMAND

CLI tool for secure file encryption and decryption.

Commands
  help    Show this help
  lock    Lock the vault
  init    Initialize the vault
  print   Print the vault
  unlock  Unlock the vault
  temp    Create a temporary vault


Options:
  -clean-print
        On print operation vault will only print the plaintext without extra info (bool, default: false)
  -do-aes256
        Use AES256 keys for asymetric vault encryption (bool, default: true) (default true)
  -do-x509
        Use X509 keys for symetric encryption (bool, default: true) (default true)
  -plain-ext string
        File extension for unencrypted plain files (string, default: txt) (default "txt")
  -private-key-path string
        Private keys path (string, default: ~/.ssh/id_rsa) (default "~/.ssh/id_rsa")
  -public-key-path string
        Public keys path (string, default: ~/.ssh/id_rsa.pub) (default "~/.ssh/id_rsa.pub")
  -vault-ext string
        File extension for encrypted vault files (string, default: vt) (default "vt")
```

# Usage
Vault operations are sub commands defined via the first command line argument.

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

# Development
You can use node.js to easily install and run nodemon or link the binary:

## Automatic building
This installs and starts a nodemon file watcher that rebuilds the binary if the sources get changed:
```sh
npm i
npm run dev
```

## Global binary linking
With global linking you can access your binary in every directory by using `vault`:
```sh
npm i -g .
```
Then use can use 'vault' everywhere:
```sh
cd ..
vault -h
```

# Contributing
Contributions to Vault are welcome!  
Interested users can refer to the guidelines provided in the [CONTRIBUTING.md](CONTRIBUTING.md) file to contribute to the project and help improve its functionality and features.

# License
Vault is licensed under the [MIT license](LICENSE), providing users with flexibility and freedom to use and modify the software according to their needs.

# Disclaimer
Vault is provided without warranties.  
Users are advised to review the accompanying license for more information on the terms of use and limitations of liability.
