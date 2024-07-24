package config

import (
	"flag"
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

func LoadConfig() AppConfig {
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

	return appConfig
}
