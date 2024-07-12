package config

import (
	"flag"
	"os"
	"strings"
)

type AppConfig struct {
	KeyDir              string
	PrivateKeyNames     []string
	PublicKeyNames      []string
	Args                []string
	VaultFileExtension  string
	PlainFileExtension  string
	BackupFileExtension string
}

func defineStringVar(
	shorthand string,
	flagName string,
	envName string,
	defaultValue string,
	description string,
) string {
	envVar, ok := os.LookupEnv(envName)
	if ok {
		defaultValue = envVar
	}

	if len(shorthand) != 0 {
		defaultValue = *flag.String(
			shorthand,
			defaultValue,
			description,
		)
	}

	return *flag.String(flagName, defaultValue, description)
}

func defineStringArrayVar(
	shorthand string,
	flagName string,
	envName string,
	defaultValue string,
	seperator string,
	description string,
) []string {
	envVar, ok := os.LookupEnv(envName)
	if ok {
		defaultValue = envVar
	}

	if len(shorthand) != 0 {
		defaultValue = *flag.String(
			shorthand,
			defaultValue,
			description,
		)
	}

	resultValue := flag.String(
		flagName,
		defaultValue,
		description+"('"+seperator+"'-seperated)",
	)

	return strings.Split(*resultValue, seperator)
}

func LoadConfig() AppConfig {
	appConfig := AppConfig{
		KeyDir: defineStringVar(
			"k",
			"key-dir",
			"VAULT_KEY_DIR",
			"~/.ssh",
			"Path to the key directory to search for asymetric keys",
		),
		PrivateKeyNames: defineStringArrayVar(
			"r",
			"private-key-names",
			"VAULT_PRIVATE_KEY_NAMES",
			"id_rsa",
			",",
			"List of private keys names",
		),
		PublicKeyNames: defineStringArrayVar(
			"u",
			"public-key-names",
			"VAULT_PUBLIC_KEY_NAMES",
			"id_rsa.pub",
			",",
			"List of public keys names",
		),
		VaultFileExtension: defineStringVar(
			"v",
			"vault-ext",
			"VAULT_EXT",
			"vt",
			"File extension for encrypted vault files",
		),
		PlainFileExtension: defineStringVar(
			"p",
			"plain-ext",
			"VAULT_PLAIN_EXT",
			"txt",
			"File extension for unencrypted plain files",
		),
		BackupFileExtension: defineStringVar(
			"b",
			"backup-ext",
			"VAULT_BACKUP_EXT",
			"bak",
			"File extension for encrypted and unencrypted backup files",
		),
		Args: []string{},
	}

	flag.Parse()

	appConfig.Args = flag.Args()

	return appConfig
}
