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

func LockOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	sourcePlainFile := targetFile + "." + appConfig.PlainFileExtension
	targetVaultFile := targetFile + "." + appConfig.VaultFileExtension

	if _, err := os.Stat(sourcePlainFile); errors.Is(err, os.ErrNotExist) {
		exitError("Plain source file '" + sourcePlainFile + "' does not exist!")
		return
	}

	plainText, err := stringfs.ReadFile(sourcePlainFile)
	if err != nil {
		exitError("Read plain source error:\n> " + err.Error())
		return
	}

	publicKey, err := cryption.LoadRsaPublicKey(appConfig.PublicKeyPath)

	if err != nil {
		exitError("Load public key error:\n> " + err.Error())
		return
	}

	lastUsedPassword, err = userin.PromptNewPassword()

	if err != nil {
		exitError("Prompt new password error:\n> " + err.Error())
		return
	}

	cipherPayload, err := VaultEncrypt(
		[]byte(plainText),
		appConfig.DoRSA,
		publicKey,
		appConfig.DoAES256,
		[]byte(lastUsedPassword),
	)

	if err != nil {
		exitError("Vault encrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFileBytes(
		targetVaultFile,
		cipherPayload,
		0640,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}

	err = stringfs.RemoveFile(sourcePlainFile)
	if err != nil {
		exitError("Remove plain source file error:\n> " + err.Error())
		return
	}

	fmt.Println("Locked!")
}
