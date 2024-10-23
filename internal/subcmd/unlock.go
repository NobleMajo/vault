package subcmd

import (
	"errors"
	"fmt"
	"os"

	"coreunit.net/vault/internal/config"
	"coreunit.net/vault/lib/stringfs"
)

func UnlockOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	sourceVaultFile := targetFile + "." + appConfig.VaultFileExtension
	targetPlainFile := targetFile + "." + appConfig.PlainFileExtension

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

	err = stringfs.SafeWriteFileBytes(
		targetPlainFile,
		plainText,
		0640,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}

	err = stringfs.RemoveFile(sourceVaultFile)
	if err != nil {
		exitError("Remove source file error:\n> " + err.Error())
		return
	}

	fmt.Println("Unlocked!")
}
