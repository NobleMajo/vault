package config

import (
	"flag"
	"os"
	"strings"
)

type AppConfig struct {
	KeyDir          string
	PrivateKeyNames []string
	PublicKeyNames  []string
	Args            []string
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

	resultValue := flag.String(flagName, defaultValue, description)
	return *resultValue
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
			"",
			"key-dir",
			"VAULT_KEY_DIR",
			"~/.ssh",
			"Path to the key directory to search for asymetric keys",
		),
		PrivateKeyNames: defineStringArrayVar(
			"",
			"private-key-names",
			"VAULT_PRIVATE_KEY_NAMES",
			"id_rsa",
			",",
			"List of private keys names",
		),
		PublicKeyNames: defineStringArrayVar(
			"",
			"public-key-names",
			"VAULT_PUBLIC_KEY_NAMES",
			"id_rsa.pub",
			",",
			"List of public keys names",
		),
		Args: []string{},
	}

	flag.Parse()

	appConfig.Args = flag.Args()

	return appConfig
}
