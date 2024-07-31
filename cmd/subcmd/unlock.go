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
