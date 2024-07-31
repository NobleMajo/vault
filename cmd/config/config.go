package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type AppConfig struct {
	Verbose             bool
	ShowVersion         bool
	ShowHelp            bool
	PrivateKeyPath      string
	PublicKeyPath       string
	Args                []string
	VaultFileExtension  string
	PlainFileExtension  string
	BackupFileExtension string
	CleanPrint          bool
	DoRSA               bool
	DoAES256            bool
	SubCommand          string
	TempDecodeSeconds   int
}

func defaultAppConfig() *AppConfig {
	return &AppConfig{
		Verbose:            false,
		ShowVersion:        false,
		ShowHelp:           false,
		PrivateKeyPath:     "~/.ssh/id_rsa",
		PublicKeyPath:      "~/.ssh/id_rsa.pub",
		Args:               []string{},
		VaultFileExtension: "vt",
		PlainFileExtension: "txt",
		CleanPrint:         false,
		DoRSA:              true,
		DoAES256:           true,
		SubCommand:         "",
		TempDecodeSeconds:  10,
	}
}

func versionCommand(appConfig *AppConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints version message",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.Args = args
			appConfig.ShowVersion = true
		},
	}

	return cmd
}

func addCryptFlags(appConfig *AppConfig, cmd *cobra.Command) {
	cmd.Flags().StringVarP(&appConfig.PrivateKeyPath, "private-key", "r", appConfig.PrivateKeyPath, "Defines the private key path (VAULT_PRIVATE_KEY_PATH)")
	cmd.Flags().StringVarP(&appConfig.PublicKeyPath, "public-key", "u", appConfig.PublicKeyPath, "Defines the public key path (VAULT_PUBLIC_KEY_PATH)")
	cmd.Flags().StringVarP(&appConfig.VaultFileExtension, "vault-ext", "e", appConfig.VaultFileExtension, "Defines the vault file extension (VAULT_EXT)")
	cmd.Flags().StringVarP(&appConfig.PlainFileExtension, "plain-ext", "p", appConfig.PlainFileExtension, "Defines the plain file extension (VAULT_PLAIN_EXT)")
	cmd.Flags().BoolVarP(&appConfig.DoRSA, "rsa", "x", appConfig.DoRSA, "Use RSA key encryption (VAULT_RSA)")
	cmd.Flags().BoolVarP(&appConfig.DoAES256, "aes", "a", appConfig.DoAES256, "Use AES256 password encryption (VAULT_AES)")
}

func lockCommand(appConfig *AppConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Locks your plain file into a vault file",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.Args = args
			appConfig.SubCommand = "lock"
		},
	}

	cmd.Aliases = append(cmd.Aliases, "loc")
	cmd.Aliases = append(cmd.Aliases, "lo")
	cmd.Aliases = append(cmd.Aliases, "l")

	addCryptFlags(appConfig, cmd)

	return cmd
}

func unlockCommand(appConfig *AppConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlock",
		Short: "Unlocks your vault file into a plain file",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.Args = args
			appConfig.SubCommand = "unlock"
		},
	}

	cmd.Aliases = append(cmd.Aliases, "unloc")
	cmd.Aliases = append(cmd.Aliases, "un")
	cmd.Aliases = append(cmd.Aliases, "u")

	addCryptFlags(appConfig, cmd)

	return cmd
}

func passwdCommand(appConfig *AppConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passwd",
		Short: "Changes the password of your vault file",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.Args = args
			appConfig.SubCommand = "unlock"
		},
	}

	cmd.Aliases = append(cmd.Aliases, "passw")
	cmd.Aliases = append(cmd.Aliases, "pass")
	cmd.Aliases = append(cmd.Aliases, "pa")

	addCryptFlags(appConfig, cmd)

	return cmd
}

func tempCommand(appConfig *AppConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "temp",
		Short: "Temporary unlocks your vault file into a plain file",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.Args = args
			appConfig.SubCommand = "temp"
		},
	}

	cmd.Aliases = append(cmd.Aliases, "tmp")
	cmd.Aliases = append(cmd.Aliases, "tm")
	cmd.Aliases = append(cmd.Aliases, "t")

	cmd.Flags().IntVarP(&appConfig.TempDecodeSeconds, "temp-seconds", "t", appConfig.TempDecodeSeconds, "Temporary decode time in seconds (VAULT_TEMP_DECODE_SECONDS)")

	addCryptFlags(appConfig, cmd)

	return cmd
}

func printCommand(appConfig *AppConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "print",
		Short: "Prints the decrypted content of your vault file",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.Args = args
			appConfig.SubCommand = "print"
		},
	}

	cmd.Aliases = append(cmd.Aliases, "prin")
	cmd.Aliases = append(cmd.Aliases, "pr")
	cmd.Aliases = append(cmd.Aliases, "p")

	cmd.Flags().BoolVarP(&appConfig.CleanPrint, "clean-print", "c", appConfig.CleanPrint, "Clean print mode (VAULT_CLEAN_PRINT)")
	addCryptFlags(appConfig, cmd)

	return cmd
}

func initCommand(appConfig *AppConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a initial encrypted vault file for default text",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.Args = args
			appConfig.SubCommand = "init"
		},
	}

	cmd.Aliases = append(cmd.Aliases, "ini")
	cmd.Aliases = append(cmd.Aliases, "in")
	cmd.Aliases = append(cmd.Aliases, "i")

	addCryptFlags(appConfig, cmd)

	return cmd
}

func loadEnvVars(appConfig *AppConfig) {
	EnvIsString("VAULT_PRIVATE_KEY_PATH", func(value string) {
		appConfig.PrivateKeyPath = value
	})

	EnvIsString("VAULT_PUBLIC_KEY_PATH", func(value string) {
		appConfig.PublicKeyPath = value
	})

	EnvIsString("VAULT_EXT", func(value string) {
		appConfig.VaultFileExtension = value
	})

	EnvIsString("VAULT_PLAIN_EXT", func(value string) {
		appConfig.PlainFileExtension = value
	})

	EnvIsBool("VAULT_RSA", func(value bool) {
		appConfig.DoRSA = value
	})

	EnvIsBool("VAULT_AES", func(value bool) {
		appConfig.DoAES256 = value
	})

	EnvIsBool("VAULT_VERBOSE", func(value bool) {
		appConfig.Verbose = value
	})

	EnvIsBool("VAULT_CLEAN_PRINT", func(value bool) {
		appConfig.CleanPrint = value
	})

	EnvIsInt("VAULT_TEMP_DECODE_SECONDS", func(value int) {
		appConfig.TempDecodeSeconds = value
	})
}

func ParseConfig(
	version string,
	commit string,
) *AppConfig {
	appConfig := defaultAppConfig()

	rootCmd := &cobra.Command{
		Use: "vault",
		Short: "File encryption and decryption cli tool written in go.\n" +
			"For more help, visit https://github.com/NobleMajo/vault",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.ShowHelp = true
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&appConfig.Verbose, "verbose", "b", appConfig.Verbose, "enable verbose mode (VAULT_VERBOSE)")
	rootCmd.Flags().BoolVarP(&appConfig.ShowVersion, "version", "v", appConfig.ShowVersion, "prints version")

	rootCmd.AddCommand(
		versionCommand(appConfig),
		initCommand(appConfig),
		printCommand(appConfig),
		lockCommand(appConfig),
		unlockCommand(appConfig),
		tempCommand(appConfig),
		passwdCommand(appConfig),
	)

	loadEnvVars(appConfig)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if appConfig.Verbose {
		fmt.Println("Verbose mode enabled")
	}

	if appConfig.ShowVersion {
		fmt.Println("Vault version " + version + ", build " + commit)
		os.Exit(0)
	}

	if appConfig.ShowHelp {
		rootCmd.Help()
		os.Exit(0)
	}

	if appConfig.SubCommand == "" {
		os.Exit(0)
	}

	return appConfig
}
