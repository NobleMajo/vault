package subcmd

import (
	"coreunit.net/vault/internal/config"
	"coreunit.net/vault/lib/stringfs"
)

func InitOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	targetVaultFile := targetFile + "." + appConfig.VaultFileExtension

	exists := stringfs.Exists(targetFile + "." + appConfig.PlainFileExtension)
	if exists {
		exitError("File plain text file '" + targetFile + "." + appConfig.PlainFileExtension + "' already exists!")
		return
	}

	exists = stringfs.Exists(targetFile + "." + appConfig.VaultFileExtension)
	if exists {
		exitError("File encrypted vault file '" + targetFile + "." + appConfig.VaultFileExtension + "' already exists!")
		return
	}

	initText := "Hello and welcome to your own vault!\n\n<3"

	loadEncryptionData(appConfig)

	cipherPayload, err := VaultEncrypt(
		[]byte(initText),
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
		targetVaultFile,
		cipherPayload,
		0644,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}
}
