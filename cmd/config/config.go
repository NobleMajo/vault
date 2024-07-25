package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type AppConfig struct {
	PrivateKeyPath      string
	PublicKeyPath       string
	Args                []string
	VaultFileExtension  string
	PlainFileExtension  string
	BackupFileExtension string
	CleanPrint          bool
	DoX509              bool
	DoAES256            bool
}

func LoadConfig(
	allowedCommands map[string]string,
) AppConfig {
	appConfig := AppConfig{
		PrivateKeyPath: DefineStringVar(
			"r",
			"private-key-path",
			"VAULT_PRIVATE_KEY_PATH",
			"~/.ssh/id_rsa",
			"Private keys path",
		),
		PublicKeyPath: DefineStringVar(
			"u",
			"public-key-path",
			"VAULT_PUBLIC_KEY_PATH",
			"~/.ssh/id_rsa.pub",
			"Public keys path",
		),
		VaultFileExtension: DefineStringVar(
			"v",
			"vault-ext",
			"VAULT_EXT",
			"vt",
			"File extension for encrypted vault files",
		),
		PlainFileExtension: DefineStringVar(
			"p",
			"plain-ext",
			"VAULT_PLAIN_EXT",
			"txt",
			"File extension for unencrypted plain files",
		),
		CleanPrint: DefineBoolVar(
			"c",
			"clean-print",
			"VAULT_CLEAN_PRINT",
			false,
			"On print operation vault will only print the plaintext without extra info",
		),
		DoX509: DefineBoolVar(
			"x",
			"do-x509",
			"VAULT_DO_X509",
			true,
			"Use X509 keys for symetric encryption",
		),
		DoAES256: DefineBoolVar(
			"a",
			"do-aes256",
			"VAULT_DO_AES256",
			true,
			"Use AES256 keys for asymetric vault encryption",
		),
		Args: []string{},
	}

	flag.Parse()

	appConfig.Args = flag.Args()

	flag.Usage = func() {
		// convert allowedCommands to a string by iterating over the keys of allowedCommands and then add them together for a good help message

		commandDescriptions := ""
		if len(allowedCommands) != 0 {
			biggestCommand := 0

			for command := range allowedCommands {
				if len(command) > biggestCommand {
					biggestCommand = len(command)
				}
			}

			if biggestCommand != 0 {
				commandDescriptions += "\nCommands\n"
				for command, description := range allowedCommands {
					commandDescriptions += "  " + command + strings.Repeat(" ", biggestCommand-len(command)) + "  " + description + "\n"
				}
			}
		}

		fmt.Printf(
			"Usage:  %s [OPTIONS] COMMAND\n"+
				"\n"+
				"CLI tool for secure file encryption and decryption.\n"+
				commandDescriptions+"\n"+
				"\n"+
				"Options:\n",
			os.Args[0],
		)
		flag.PrintDefaults()
		fmt.Println()

	}

	return appConfig
}
