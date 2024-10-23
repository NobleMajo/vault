package subcmd

import (
	"errors"
	"fmt"
	"os"

	"coreunit.net/vault/internal/config"
	"coreunit.net/vault/lib/stringfs"
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

	loadDecryptionData(appConfig)

	plainText, err := VaultDecrypt(
		[]byte(vaultRaw),
		!appConfig.DisableRSA,
		lastUsedPrivateKey,
		!appConfig.DisableAES256,
		[]byte(lastUsedPassword),
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
