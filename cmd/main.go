package main

import (
	"bufio"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"coreunit.net/vault/cmd/config"
	"coreunit.net/vault/internal/cryption"
	"coreunit.net/vault/internal/stringfs"
	"golang.org/x/term"
)

func main() {
	appConfig := config.ParseConfig()

	stringfs.ParsePath(&appConfig.PublicKeyPath)
	stringfs.ParsePath(&appConfig.PrivateKeyPath)
	targetFile := targetFile(appConfig)

	if appConfig.SubCommand == "lock" {
		lockOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "init" {
		initOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "print" {
		printOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "unlock" {
		unlockOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "temp" {
		tempOperation(
			targetFile,
			appConfig,
		)
	} else {
		fmt.Fprintf(
			os.Stderr,
			"%s: '"+appConfig.SubCommand+"' is not a command.\n"+
				"See '%s help'",
			os.Args[0],
			os.Args[0],
		)
	}
}

func targetFile(
	appConfig *config.AppConfig,
) string {
	if len(appConfig.Args) >= 2 {
		targetFile := appConfig.Args[1]

		if strings.HasSuffix(targetFile, "."+appConfig.VaultFileExtension) {
			targetFile = targetFile[:len(targetFile)-len(appConfig.VaultFileExtension)-1]
		} else if strings.HasSuffix(targetFile, "."+appConfig.PlainFileExtension) {
			targetFile = targetFile[:len(targetFile)-len(appConfig.PlainFileExtension)-1]
		}

		return targetFile
	}

	return "vault"
}

func exitError(message string) {
	fmt.Fprintln(
		os.Stderr,
		message,
	)
	os.Exit(1)
}

func PromptNewPassword() (string, error) {
	fmt.Println("Enter your new vault password:")
	newPassword, err := ReadPassword()
	if err != nil {
		return "", err
	}
	fmt.Println("Re-enter your new vault password:")

	newPassword2, err := ReadPassword()
	if err != nil {
		return "", err
	}

	if newPassword != newPassword2 {
		return "", fmt.Errorf("passwords don't match")
	}

	return newPassword, nil
}

func PromptPassword() (string, error) {
	fmt.Println("Enter your vault password:")
	newPassword, err := ReadPassword()
	if err != nil {
		return "", err
	}

	return newPassword, nil
}

func VaultEncrypt(
	payload []byte,
	doX509 bool,
	X509PublicKey *rsa.PublicKey,
	doAES256 bool,
	AES256Key []byte,
) ([]byte, error) {
	if !doX509 && !doAES256 {
		return nil, fmt.Errorf("no encryption method selected")
	}
	var err error

	if doAES256 {
		payload, err = cryption.AES256Encrypt(AES256Key, payload)
		if err != nil {
			return nil, fmt.Errorf("AES256 encrypt error, maybe wrong password:\n> %v", err)
		}
	}

	if doX509 {
		payload, err = cryption.X509AES256Encrypt(X509PublicKey, payload)
		if err != nil {
			return nil, fmt.Errorf("x509 encrypt error:\n> %v", err)
		}
	}

	return payload, nil
}

func VaultDecrypt(
	payload []byte,
	doX509 bool,
	X509PrivateKey *rsa.PrivateKey,
	doAES256 bool,
	AES256Key []byte,
) ([]byte, error) {
	if !doX509 && !doAES256 {
		return nil, fmt.Errorf("no decryption method selected")
	}
	var err error

	if doX509 {
		payload, err = cryption.X509AES256Decrypt(X509PrivateKey, payload)
		if err != nil {
			return nil, fmt.Errorf("x509 decrypt error:\n> %v", err)
		}
	}

	if doAES256 {
		payload, err = cryption.AES256Decrypt(AES256Key, payload)
		if err != nil {
			return nil, fmt.Errorf("AES256 decrypt error, maybe wrong password:\n> %v", err)
		}
	}

	return payload, nil
}

func ReadPassword() (string, error) {
	rawData, err := term.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func ReadLine() (string, error) {
	fmt.Print("> ")
	buf := bufio.NewReader(os.Stdin)
	rawData, err := buf.ReadBytes('\n')

	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func initOperation(
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

	publicKey, err := cryption.LoadRsaPublicKey(appConfig.PublicKeyPath)

	if err != nil {
		exitError("Load public key error:\n> " + err.Error())
		return
	}

	password, err := PromptNewPassword()

	if err != nil {
		exitError("Prompt new password error:\n> " + err.Error())
		return
	}

	cipherPayload, err := VaultEncrypt(
		[]byte(initText),
		appConfig.DoRSA,
		publicKey,
		appConfig.DoAES256,
		[]byte(password),
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

var lastUsedPassword string

func lockOperation(
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

	lastUsedPassword, err = PromptNewPassword()

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

	fmt.Println("Locked!")
}

func unlockOperation(
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

	password, err := PromptPassword()

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

func tempOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	var err error
	decodeTime := 10

	if len(appConfig.Args) == 3 && len(appConfig.Args[2]) != 0 {
		rawDecodeTime := appConfig.Args[2]

		decodeTime, err = strconv.Atoi(strings.TrimSpace(rawDecodeTime))

		if err != nil {
			exitError("Error parsing decode time to number '" + rawDecodeTime + "':\n> " + err.Error())
			return
		}
	}

	unlockOperation(targetFile, appConfig)

	sourceVaultFile := targetFile + "." + appConfig.VaultFileExtension
	targetPlainFile := targetFile + "." + appConfig.PlainFileExtension

	fmt.Println("Unlocked for " + strconv.Itoa(decodeTime) + " seconds!")
	time.Sleep(time.Duration(decodeTime) * time.Second)
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

	fmt.Println("Locked again!")
}

func printOperation(
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

	password, err := PromptPassword()

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
