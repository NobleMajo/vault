# Table of Contents
- [Table of Contents](#table-of-contents)
- [About](#about)
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

# About
Vault is a small and simple CLI tool that encrypt and decrypt plain `.txt` files into vault-files (`.vt`, a custom file format).

The idea behind this tool is to have a CLI utility that can quickly and easily encrypt individual files, allowing users to securely store API tokens, secrets, credentials, or any private data on their own disk.

## Advertising
*Are you also a small software developer or admin with lots of API keys, encryption keys or other secrets and credentials?*
*Or do you simply have logs or plain text files that you want to send to someone securely?*
**Then I have exactly what you are looking for today!**

*Hold on tight and take a closer look at this command line interface tool because it might meet your exact needs.*

## Encryption
Vault uses asymmetric RSA encryption and symmetric AES-256 encryption to keep your data as secure as possible.
To do this, vault uses private and public key on disk (default: `~/.ssh/id_rsa.pub`) and also asks you for a password.

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
Build the vault binary from source code:
```sh
go build -o bin/vault cmd/main.go
# or npm run build
```

## Help
Run the help command on the binary
```sh
./vault -h
```

Note: Almost everything in vault can be set using a flag and the corresponding environment variable.
Check out the help page:
```rust
Usage of vault:
  -b string
        File extension for encrypted and unencrypted backup files (default "bak")
  -backup-ext string
        File extension for encrypted and unencrypted backup files (default "bak")
  -k string
        Path to the key directory to search for asymetric keys (default "~/.ssh")
  -key-dir string
        Path to the key directory to search for asymetric keys (default "~/.ssh")
  -p string
        File extension for unencrypted plain files (default "txt")
  -plain-ext string
        File extension for unencrypted plain files (default "txt")
  -private-key-names string
        List of private keys names(','-seperated) (default "id_rsa")
  -public-key-names string
        List of public keys names(','-seperated) (default "id_rsa.pub")
  -r string
        List of private keys names (default "id_rsa")
  -u string
        List of public keys names (default "id_rsa.pub")
  -v string
        File extension for encrypted vault files (default "vt")
  -vault-ext string
        File extension for encrypted vault files (default "vt")
```

# Usage
Vault operations are sub commands defined via the first command line argument.

### lock
Add some content to your `vault.txt` and lock it:
```sh
vim vault.txt
vault lock
```

**OR**

### init
Create a new locked vault file:
```sh
vault init
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
