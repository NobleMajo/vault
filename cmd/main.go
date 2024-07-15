package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"syscall"
	"time"

	"coreunit.net/vault/internal/stringcrypt"
	"coreunit.net/vault/internal/stringfs"
	"golang.org/x/term"
)

func main() {
	appConfig := LoadConfig()

	err := stringfs.ParsePath(&appConfig.KeyDir)
	if err != nil {
		exitError("Key dir '" + appConfig.KeyDir + "' is not a valid path:\n> " + err.Error())
	}

	allowedOperations := []string{"lock", "unlock", "temp", "init", "print", "help"}

	if len(appConfig.Args) == 0 {
		exitError("First argument needs to be an operation: '" + strings.Join(allowedOperations, "', '") + "'!")
		return
	}

	rawOperation := strings.ToLower(appConfig.Args[0])

	if rawOperation == "help" {
		fmt.Fprintf(os.Stderr, "Help: Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
		return
	}

	if !slices.Contains(allowedOperations, rawOperation) {
		exitError(
			"First argument needs to be an operation: '" + strings.Join(allowedOperations, "', '") + "'\n" +
				"Not: '" + rawOperation + "'",
		)
		return
	}

	var targetFile string

	if len(appConfig.Args) >= 2 {
		targetFile = appConfig.Args[1]

		if strings.HasSuffix(targetFile, "."+appConfig.VaultFileExtension) {
			targetFile = targetFile[:len(targetFile)-len(appConfig.VaultFileExtension)-1]
		} else if strings.HasSuffix(targetFile, "."+appConfig.PlainFileExtension) {
			targetFile = targetFile[:len(targetFile)-len(appConfig.PlainFileExtension)-1]
		}
	} else {
		targetFile = "vault"
	}

	if rawOperation == "lock" {
		lockOperation(
			targetFile,
			appConfig,
		)
		return
	} else if rawOperation == "init" {
		initOperation(
			targetFile,
			appConfig,
		)
		return
	} else if rawOperation == "print" {
		printOperation(
			targetFile,
			appConfig,
		)
		return
	} else if rawOperation == "unlock" || rawOperation == "temp" {
		unlockOperation(
			targetFile,
			rawOperation == "temp",
			appConfig,
		)
		return
	}

	exitError("Unknown and not implemented operation: '" + rawOperation + "'")
}

func exitError(message string) {
	fmt.Fprintf(os.Stderr, "Help: Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Println("ERROR: " + message)
	os.Exit(1)
}

func initOperation(
	targetFile string,
	appConfig AppConfig,
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

	publicKey, err := stringcrypt.LoadPublicKey(
		appConfig.KeyDir,
		appConfig.PublicKeyNames,
	)
	if err != nil {
		exitError("Cant load public key:\n> " + err.Error())
		return
	}

	fmt.Println("Enter your new password:")
	password, err := ReadPassword()

	if err != nil {
		exitError("Password input error:\n> " + err.Error())
		return
	}

	initText := "Hello and welcome to your own vault!\n\n<3"
	encodedText, err := VaultEncrypt(initText, publicKey, password)

	if err != nil {
		exitError("Vault encrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFile(
		targetVaultFile,
		"."+appConfig.BackupFileExtension,
		string(encodedText),
		0644,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}
}

func lockOperation(
	targetFile string,
	appConfig AppConfig,
) {
	// ENCODING
	sourcePlainFile := targetFile + "." + appConfig.PlainFileExtension
	targetVaultFile := targetFile + "." + appConfig.VaultFileExtension

	if _, err := os.Stat(sourcePlainFile); errors.Is(err, os.ErrNotExist) {
		exitError("Plain source file '" + sourcePlainFile + "' does not exist!")
		return
	}

	plainText, err := stringfs.SafeReadFile(sourcePlainFile, "."+appConfig.BackupFileExtension)
	if err != nil {
		exitError("Read plain source error:\n> " + err.Error())
		return
	}

	publicKey, err := stringcrypt.LoadPublicKey(
		appConfig.KeyDir,
		appConfig.PublicKeyNames,
	)
	if err != nil {
		exitError("Cant load public key:\n> " + err.Error())
		return
	}

	fmt.Println("Enter your new password:")
	password, err := ReadPassword()

	if err != nil {
		exitError("Password input error:\n> " + err.Error())
		return
	}

	encodedText, err := VaultEncrypt(plainText, publicKey, password)

	if err != nil {
		exitError("Vault encrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFile(
		targetVaultFile,
		"."+appConfig.BackupFileExtension,
		string(encodedText),
		0644,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeRemoveFile(sourcePlainFile, "."+appConfig.BackupFileExtension)
	if err != nil {
		exitError("Remove source file error:\n> " + err.Error())
		return
	}

	fmt.Println("Locked!")
	os.Exit(0)
}

func unlockOperation(
	targetFile string,
	isTemp bool,
	appConfig AppConfig,
) {
	sourceVaultFile := targetFile + "." + appConfig.VaultFileExtension
	targetPlainFile := targetFile + "." + appConfig.PlainFileExtension

	if _, err := os.Stat(sourceVaultFile); errors.Is(err, os.ErrNotExist) {
		exitError("Source vault file '" + sourceVaultFile + "' does not exist!")
		return
	}

	vaultRaw, err1, vaultBackupRaw, err2 := stringfs.SafeReadBothFiles(sourceVaultFile, "."+appConfig.BackupFileExtension)

	if err1 != nil && err2 != nil {
		exitError("Read plain source error:\n" + err1.Error())
		return
	}

	fmt.Println("Enter your password:")
	password, err := ReadPassword()

	if err != nil {
		exitError("Password input error:\n> " + err.Error())
		return
	}

	privateKey, err := stringcrypt.LoadPrivateKey(
		appConfig.KeyDir,
		appConfig.PrivateKeyNames,
		func() (string, error) {
			fmt.Println("Enter your private key passphrase:")
			password, err := ReadPassword()
			if err != nil {
				return "", err
			}
			return password, nil
		},
	)
	if err != nil {
		exitError("Cant load public key:\n> " + err.Error())
		return
	}

	decodedText, err := VaultDecrypt(vaultRaw, privateKey, password)
	if err != nil {
		fmt.Println("The original vault file has an X509 decryption error, try using backup file instead")
		decodedText, err = VaultDecrypt(vaultBackupRaw, privateKey, password)
		if err != nil {
			exitError("Decrypt error:\n> " + err.Error())
			return
		}
	}

	err = stringfs.SafeWriteFile(
		targetPlainFile,
		"."+appConfig.BackupFileExtension,
		decodedText,
		0644,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeRemoveFile(sourceVaultFile, "."+appConfig.BackupFileExtension)
	if err != nil {
		exitError("Remove source file error:\n> " + err.Error())
		return
	}

	if isTemp {
		tempOperation(
			targetPlainFile,
			password,
			sourceVaultFile,
			appConfig,
		)
		return
	}

	fmt.Println("Unlocked!")
	os.Exit(0)
}

func printOperation(
	targetFile string,
	appConfig AppConfig,
) {
	sourceVaultFile := targetFile + "." + appConfig.VaultFileExtension

	if _, err := os.Stat(sourceVaultFile); errors.Is(err, os.ErrNotExist) {
		exitError("Source vault file '" + sourceVaultFile + "' does not exist!")
		return
	}

	vaultRaw, err1, vaultBackupRaw, err2 := stringfs.SafeReadBothFiles(sourceVaultFile, "."+appConfig.BackupFileExtension)

	if err1 != nil && err2 != nil {
		exitError("Read plain source error:\n" + err1.Error())
		return
	}

	fmt.Println("Enter your password:")
	password, err := ReadPassword()

	if err != nil {
		exitError("Password input error:\n> " + err.Error())
		return
	}

	privateKey, err := stringcrypt.LoadPrivateKey(
		appConfig.KeyDir,
		appConfig.PrivateKeyNames,
		func() (string, error) {
			fmt.Println("Enter your private key passphrase:")
			password, err := ReadPassword()
			if err != nil {
				return "", err
			}
			return password, nil
		},
	)
	if err != nil {
		exitError("Cant load public key:\n> " + err.Error())
		return
	}

	decodedText, err := VaultDecrypt(vaultRaw, privateKey, password)
	if err != nil {
		fmt.Println("The original vault file has an X509 decryption error, try using backup file instead")
		decodedText, err = VaultDecrypt(vaultBackupRaw, privateKey, password)
		if err != nil {
			exitError("Decrypt error:\n> " + err.Error())
			return
		}
	}

	var result string = ""
	width, height, err := term.GetSize(0)

	firstMessage := "Vault Content:"
	lastMessage := "Don't forget to use 'clear'!"

	if err != nil || width < 16 {
		result += firstMessage + "\n"
		if err != nil || height > 16 {
			result += strings.Repeat("-", width) + "\n"
		}
		result += "\n"
		result += decodedText + "\n"
		result += "\n"
		if err != nil || height > 16 {
			result += strings.Repeat("-", width) + "\n"
		}
		result += lastMessage
	} else {
		firstSpaces := width/2 - (len(firstMessage) / 2)
		lastSpaces := width/2 - (len(lastMessage) / 2)

		if height > 16 {
			result += "\n"
		}
		result += strings.Repeat(" ", firstSpaces) + firstMessage + "\n"
		if height > 16 {
			result += "\n"
		}
		result += "#" + strings.Repeat("-", width-2) + "#\n"
		if height > 16 {
			result += "|\n"
		}
		result += "|  " + strings.Join(strings.Split(decodedText, "\n"), "\n|  ") + "\n"
		if height > 16 {
			result += "|\n"
		}
		result += "#" + strings.Repeat("-", width-2) + "#\n"
		if height > 16 {
			result += "\n"
		}
		result += strings.Repeat(" ", lastSpaces) + lastMessage
	}

	fmt.Println(result)
	os.Exit(0)
}

func tempOperation(
	targetPlainFile string,
	password string,
	sourceVaultFile string,
	appConfig AppConfig,
) {

	fmt.Println("Unlocked for 5 seconds!")
	time.Sleep(5 * time.Second)
	fmt.Println("Lock vault now...")

	plainText, err := stringfs.SafeReadFile(targetPlainFile, "."+appConfig.BackupFileExtension)
	if err != nil {
		exitError("Read plain source error:\n> " + err.Error())
		return
	}

	publicKey, err := stringcrypt.LoadPublicKey(
		appConfig.KeyDir,
		appConfig.PublicKeyNames,
	)
	if err != nil {
		exitError("Cant load public key:\n> " + err.Error())
		return
	}

	encodedText, err := VaultEncrypt(plainText, publicKey, password)

	if err != nil {
		exitError("Vault encrypt error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeWriteFile(
		sourceVaultFile,
		"."+appConfig.BackupFileExtension,
		string(encodedText),
		0644,
	)

	if err != nil {
		exitError("Write file error:\n> " + err.Error())
		return
	}

	err = stringfs.SafeRemoveFile(targetPlainFile, "."+appConfig.BackupFileExtension)
	if err != nil {

		exitError("Remove source file error:\n> " + err.Error())
		return
	}

	fmt.Println("Locked again!")
	os.Exit(0)
}

func VaultEncrypt(
	rawContent string,
	publicKey string,
	password string,
) (string, error) {
	encodedContent, err := stringcrypt.AesEncrypt(password, rawContent)
	if err != nil {
		return "", fmt.Errorf("aes encrypt error:\n> %v", err)
	}

	encodedContent, err = stringcrypt.X509Encrypt(publicKey, encodedContent)
	if err != nil {
		return "", fmt.Errorf("x509 encrypt error:\n> %v", err)
	}

	return encodedContent, nil
}

func VaultDecrypt(
	encodedContent string,
	privateKey string,
	password string,
) (string, error) {
	decodedText, err := stringcrypt.X509Decrypt(privateKey, encodedContent)
	if err != nil {
		return "", fmt.Errorf("x509 decrypt error:\n> %v", err)
	}

	decodedText, err = stringcrypt.AesDecrypt(password, decodedText)
	if err != nil {
		return "", fmt.Errorf("aes decrypt error, maybe wrong password:\n> %v", err)
	}

	return decodedText, nil
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
