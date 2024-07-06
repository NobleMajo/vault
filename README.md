# Table of Contents
- [Table of Contents](#table-of-contents)
- [About](#about)
- [Future](#future)
- [Getting Started](#getting-started)
  - [Requirements](#requirements)
  - [Build](#build)
  - [Run](#run)
- [Contributing](#contributing)
- [License](#license)
- [Disclaimer](#disclaimer)

# About
Vault is a small and simple CLI tool that encodes plain `.txt` files into encrypted `.vt` files.
`.vt` is a custom file format for vault.

The idea behind this tool is to have a CLI utility that can quickly and easily encrypt individual files, allowing users to securely store API tokens, secrets, credentials, or any private data on their own disk.

# Future
Currently, the Vault file is encrypted using AES-256, but in the future it will also use the private and public keys from the .ssh directory to add a addiitonal asymetic encryption by default.
This can be turned off via a flag.

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
```sh
go build -o bin/vault cmd/main.go
# or npm run build
```

## Run
Add some content to your `vault.txt` and lock it:
```sh
vim vault.txt
./vault lock
```

Unlock the file for 30 seconds, changes get saved:
```sh
./vault temp
```

Unlock forever:
```sh
./vault unlock
```

To choose a other file then the `vault.txt` use the second argument without extentions:
 ```sh
./vault lock <filename>
./vault temp <filename>
./vault unlock <filename>
```

# Contributing
Contributions to Vault are welcome!  
Interested users can refer to the guidelines provided in the [CONTRIBUTING.md](CONTRIBUTING.md) file to contribute to the project and help improve its functionality and features.

# License
Vault is licensed under the [MIT license](LICENSE), providing users with flexibility and freedom to use and modify the software according to their needs.

# Disclaimer
Vault is provided without warranties.  
Users are advised to review the accompanying license for more information on the terms of use and limitations of liability.
