package subcmd

import (
	"errors"
	"fmt"
	"os"

	"coreunit.net/vault/internal/config"
	"coreunit.net/vault/lib/stringfs"
)

func PasswdOperation(
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
		exitError("Vault decrypt error:\n> " + err.Error())
		return
	}

	lastUsedPassword = ""
	loadEncryptionData(appConfig)

	cipherPayload, err := VaultEncrypt(
		[]byte(plainText),
		!appConfig.DisableRSA,
		lastUsedPublicKey,
		!appConfig.DisableAES256,
		[]byte(lastUsedPassword),
	)

	if err != nil {
		exitError("Vault encrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFileBytes(
		sourceVaultFile,
		cipherPayload,
		0640,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}

	fmt.Println("Password changed!")
}
