package subcmd

import (
	"errors"
	"fmt"
	"os"

	"coreunit.net/vault/cmd/config"
	"coreunit.net/vault/internal/cryption"
	"coreunit.net/vault/internal/stringfs"
	"coreunit.net/vault/internal/userin"
)

func PrintOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	sourceVaultFile := targetFile + "." + appConfig.VaultFileExtension

	if _, err := os.Stat(sourceVaultFile); errors.Is(err, os.ErrNotExist) {
		exitError("Source vault file '" + sourceVaultFile + "' does not exist!")
		return
	}

	vaultRaw, err := stringfs.ReadFile(sourceVaultFile)

	if err != nil {
		exitError("Error while read vault source from '" + sourceVaultFile + "':\n> " + err.Error())
		return
	}

	privateKey, err := cryption.LoadRsaPrivateKey(appConfig.PrivateKeyPath)

	if err != nil {
		exitError("Load private key error:\n> " + err.Error())
		return
	}

	password, err := userin.PromptPassword()

	if err != nil {
		exitError("Prompt password error:\n> " + err.Error())
		return
	}

	plainText, err := VaultDecrypt(
		[]byte(vaultRaw),
		appConfig.DoRSA,
		privateKey,
		appConfig.DoAES256,
		[]byte(password),
	)

	if err != nil {
		exitError("Decrypt error:\n> " + err.Error())
		return
	}

	if appConfig.CleanPrint {
		fmt.Println(plainText)
	} else {
		fmt.Println(
			"### Vault Content:\n\n" +
				string(plainText) + "\n\n" +
				"### Don't forget to clear!",
		)
	}
}
