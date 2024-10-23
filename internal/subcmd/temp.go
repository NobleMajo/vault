package subcmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"coreunit.net/vault/internal/config"
	"coreunit.net/vault/lib/stringfs"
)

func TempOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	sourceVaultFile := targetFile + "." + appConfig.VaultFileExtension
	targetPlainFile := targetFile + "." + appConfig.PlainFileExtension

	fmt.Println(
		"Temporary unlock vault for " +
			strconv.Itoa(appConfig.TempDecodeSeconds) +
			" seconds...",
	)

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

	decryptedPlainText, err := VaultDecrypt(
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
		decryptedPlainText,
		0640,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}

	fmt.Println("Unlocked! Wait for " + strconv.Itoa(appConfig.TempDecodeSeconds) + " seconds...")
	time.Sleep(time.Duration(appConfig.TempDecodeSeconds) * time.Second)
	fmt.Println("Lock vault now again...")

	if _, err := os.Stat(targetPlainFile); errors.Is(err, os.ErrNotExist) {
		exitError("Plain source file '" + targetPlainFile + "' does not exist!")
		return
	}

	plainText, err := stringfs.ReadFile(targetPlainFile)
	if err != nil {
		exitError("Read plain source error:\n> " + err.Error())
		return
	}

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

	err = stringfs.RemoveFile(targetPlainFile)
	if err != nil {
		exitError("Remove plain source file error:\n> " + err.Error())
		return
	}

	fmt.Println("Locked again!")
}
