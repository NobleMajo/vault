package subcmd

import (
	"fmt"
	"strconv"
	"time"

	"coreunit.net/vault/cmd/config"
	"coreunit.net/vault/internal/cryption"
	"coreunit.net/vault/internal/stringfs"
)

var lastUsedPassword string

func TempOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	var err error

	UnlockOperation(targetFile, appConfig)

	sourceVaultFile := targetFile + "." + appConfig.VaultFileExtension
	targetPlainFile := targetFile + "." + appConfig.PlainFileExtension

	fmt.Println("Unlocked for " + strconv.Itoa(appConfig.TempDecodeSeconds) + " seconds!")
	time.Sleep(time.Duration(appConfig.TempDecodeSeconds) * time.Second)
	fmt.Println("Lock vault now...")

	plainText, err := stringfs.ReadFile(targetPlainFile)
	if err != nil {
		exitError("Read plain source error:\n> " + err.Error())
		return
	}

	publicKey, err := cryption.LoadRsaPublicKey(appConfig.PublicKeyPath)

	if err != nil {
		exitError("Load public key error:\n> " + err.Error())
		return
	}

	encodedText, err := VaultEncrypt(
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
		sourceVaultFile,
		encodedText,
		0644,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}

	err = stringfs.RemoveFile(targetPlainFile)
	if err != nil {
		exitError("Remove temp plain file error:\n> " + err.Error())
		return
	}

	fmt.Println("Locked again!")
}
