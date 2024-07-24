package main

import (
	"bufio"
	"crypto/rsa"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"coreunit.net/vault/cmd/config"
	"coreunit.net/vault/internal/stringcrypt"
	"coreunit.net/vault/internal/stringfs"
	"golang.org/x/term"
)

var commands []string

func main() {
	appConfig := config.LoadConfig()
	commands = []string{"help", "lock", "init", "print", "unlock", "temp"}

	if len(appConfig.Args) == 0 {
		PrintHelp()
		return
	}

	rawOperation := strings.ToLower(appConfig.Args[0])

	if rawOperation == "help" {
		PrintHelp()
		return
	}

	stringfs.ParsePath(&appConfig.PublicKeyPath)
	stringfs.ParsePath(&appConfig.PrivateKeyPath)
	targetFile := targetFile(appConfig)

	if rawOperation == "lock" {
		lockOperation(
			targetFile,
			appConfig,
		)
	} else if rawOperation == "init" {
		initOperation(
			targetFile,
			appConfig,
		)
	} else if rawOperation == "print" {
		printOperation(
			targetFile,
			appConfig,
		)
	} else if rawOperation == "unlock" {
		unlockOperation(
			targetFile,
			appConfig,
		)
	} else if rawOperation == "temp" {
		tempOperation(
			targetFile,
			appConfig,
		)
	} else {
		fmt.Fprintf(
			os.Stderr,
			"%s: '"+rawOperation+"' is not a command.\n"+
				"See '%s help'",
			os.Args[0],
			os.Args[0],
		)
	}
}

func targetFile(
	appConfig config.AppConfig,
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

func initOperation(
	targetFile string,
	appConfig config.AppConfig,
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

	publicKey, err := stringcrypt.LoadRsaPublicKey(appConfig.PublicKeyPath)

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
		initText,
		appConfig.DoX509,
		publicKey,
		appConfig.DoAES256,
		password,
	)

	if err != nil {
		exitError("Vault encrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFile(
		targetVaultFile,
		cipherPayload,
		0644,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}
}

func lockOperation(
	targetFile string,
	appConfig config.AppConfig,
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

	publicKey, err := stringcrypt.LoadRsaPublicKey(appConfig.PublicKeyPath)

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
		plainText,
		appConfig.DoX509,
		publicKey,
		appConfig.DoAES256,
		password,
	)

	if err != nil {
		exitError("Vault encrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFile(
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
	appConfig config.AppConfig,
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

	privateKey, err := stringcrypt.LoadRsaPrivateKey(appConfig.PrivateKeyPath)

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
		vaultRaw,
		appConfig.DoX509,
		privateKey,
		appConfig.DoAES256,
		password,
	)

	if err != nil {
		exitError("Vault decrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFile(
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
	appConfig config.AppConfig,
) {
	unlockOperation(targetFile, appConfig)

	sourceVaultFile := targetFile + "." + appConfig.VaultFileExtension
	targetPlainFile := targetFile + "." + appConfig.PlainFileExtension

	fmt.Println("Unlocked for 5 seconds!")
	time.Sleep(5 * time.Second)
	fmt.Println("Lock vault now...")

	plainText, err := stringfs.ReadFile(targetPlainFile)
	if err != nil {
		exitError("Read plain source error:\n> " + err.Error())
		return
	}

	publicKey, err := stringcrypt.LoadRsaPublicKey(appConfig.PublicKeyPath)

	if err != nil {
		exitError("Load public key error:\n> " + err.Error())
		return
	}

	password, err := PromptNewPassword()

	if err != nil {
		exitError("Prompt new password error:\n> " + err.Error())
		return
	}

	encodedText, err := VaultEncrypt(
		plainText,
		appConfig.DoX509,
		publicKey,
		appConfig.DoAES256,
		password,
	)

	if err != nil {
		exitError("Vault encrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFile(
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
	appConfig config.AppConfig,
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

	privateKey, err := stringcrypt.LoadRsaPrivateKey(appConfig.PrivateKeyPath)

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
		vaultRaw,
		appConfig.DoX509,
		privateKey,
		appConfig.DoAES256,
		password,
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
				plainText + "\n\n" +
				"### Don't forget to clear!",
		)
	}
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
	payload string,
	doX509 bool,
	X509PublicKey *rsa.PublicKey,
	doAES256 bool,
	AES256Key string,
) (string, error) {
	if !doX509 && !doAES256 {
		return "", fmt.Errorf("no encryption method selected")
	}
	var err error

	if doAES256 {
		payload, err = stringcrypt.AES256Encrypt(AES256Key, payload)
		if err != nil {
			return "", fmt.Errorf("AES256 encrypt error, maybe wrong password:\n> %v", err)
		}
	}

	if doX509 {
		payload, err = stringcrypt.X509Encrypt(X509PublicKey, payload)
		if err != nil {
			return "", fmt.Errorf("x509 encrypt error:\n> %v", err)
		}
	}

	return payload, nil
}

func VaultDecrypt(
	payload string,
	doX509 bool,
	X509PrivateKey *rsa.PrivateKey,
	doAES256 bool,
	AES256Key string,
) (string, error) {
	if !doX509 && !doAES256 {
		return "", fmt.Errorf("no decryption method selected")
	}
	var err error

	if doX509 {
		payload, err = stringcrypt.X509Decrypt(X509PrivateKey, payload)
		if err != nil {
			return "", fmt.Errorf("x509 decrypt error:\n> %v", err)
		}
	}

	if doAES256 {
		payload, err = stringcrypt.AES256Decrypt(AES256Key, payload)
		if err != nil {
			return "", fmt.Errorf("AES256 decrypt error, maybe wrong password:\n> %v", err)
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

func PrintHelp() {
	fmt.Printf(
		"Usage:  %s [OPTIONS] COMMAND\n"+
			"\n"+
			"CLI tool for secure file encryption and decryption.\n"+
			"\n"+
			"Commands:\n"+
			"  "+strings.Join(
			commands,
			",\n  ",
		)+"\n"+
			"\n"+
			"Options:\n",
		os.Args[0],
	)
	flag.PrintDefaults()
	fmt.Println()
}
