package main

import (
	"bufio"
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

func main() {
	appConfig := config.LoadConfig()

	if len(appConfig.Args) < 1 {
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("ERROR: First argument needs to be an operation: 'unlock', 'lock', 'temp', 'help'!")
		os.Exit(1)
		return
	}

	rawOperation := strings.ToLower(appConfig.Args[0])

	if rawOperation == "help" {
		fmt.Fprintf(os.Stderr, "Help: Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
		return
	}

	if rawOperation != "lock" && rawOperation != "unlock" && rawOperation != "temp" {
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("ERROR: First argument needs to be an operation: 'unlock', 'lock', 'temp'\nNot: '" + rawOperation + "'")
		os.Exit(1)
		return
	}

	vaultFileExtension := "vt"
	plainFileExtension := "txt"
	backupFileExtension := "bak"
	var targetFile string

	if len(appConfig.Args) >= 2 {
		targetFile = appConfig.Args[1]

		if strings.HasSuffix(targetFile, "."+vaultFileExtension) {
			targetFile = targetFile[:len(targetFile)-len(vaultFileExtension)-1]
		} else if strings.HasSuffix(targetFile, "."+plainFileExtension) {
			targetFile = targetFile[:len(targetFile)-len(plainFileExtension)-1]
		}
	} else {
		targetFile = "vault"
	}

	if rawOperation == "lock" {
		// ENCODING
		sourcePlainFile := targetFile + "." + plainFileExtension
		targetVaultFile := targetFile + "." + vaultFileExtension

		if _, err := os.Stat(sourcePlainFile); errors.Is(err, os.ErrNotExist) {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("ERROR: Plain source file '" + sourcePlainFile + "' does not exist!")
			os.Exit(1)
			return
		}

		plainText, err := stringfs.SafeReadFile(sourcePlainFile, "."+backupFileExtension)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Read plain source error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		fmt.Println("Enter your new password:")
		password, err := ReadPassword()

		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Password input error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		encodedText, err := stringcrypt.AesEncrypt(password, plainText)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Encrypt error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		err = stringfs.SafeWriteFile(
			targetVaultFile,
			"."+backupFileExtension,
			string(encodedText),
			0644,
		)

		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Write file error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		err = stringfs.SafeRemoveFile(sourcePlainFile, "."+backupFileExtension)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Remove source file error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		fmt.Println("Locked!")
		os.Exit(0)
		return
	} else if rawOperation == "unlock" || rawOperation == "temp" {
		// DECODING
		sourceVaultFile := targetFile + "." + vaultFileExtension
		targetPlainFile := targetFile + "." + plainFileExtension

		if _, err := os.Stat(sourceVaultFile); errors.Is(err, os.ErrNotExist) {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("ERROR: Source vault file '" + sourceVaultFile + "' does not exist!")
			os.Exit(1)
			return
		}

		vaultRaw, err := stringfs.SafeReadFile(sourceVaultFile, "."+backupFileExtension)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Read plain source error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		fmt.Println("Enter your password:")
		password, err := ReadPassword()

		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Password input error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		encodedText, err := stringcrypt.AesDecrypt(password, vaultRaw)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Encrypt error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		err = stringfs.SafeWriteFile(
			targetPlainFile,
			"."+backupFileExtension,
			encodedText,
			0644,
		)

		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Write file error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		err = stringfs.SafeRemoveFile(sourceVaultFile, "."+backupFileExtension)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Remove source file error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		if rawOperation == "unlock" {
			fmt.Println("Unlocked!")
			os.Exit(0)
			return
		}

		fmt.Println("Unlocked for 30 seconds!")
		time.Sleep(30 * time.Second)
		fmt.Println("Lock vault now...")

		plainText, err := stringfs.SafeReadFile(targetPlainFile, "."+backupFileExtension)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Read plain source error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		encodedText, err = stringcrypt.AesEncrypt(password, plainText)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Encrypt error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		err = stringfs.SafeWriteFile(
			sourceVaultFile,
			"."+backupFileExtension,
			string(encodedText),
			0644,
		)

		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Write file error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		err = stringfs.SafeRemoveFile(targetPlainFile, "."+backupFileExtension)
		if err != nil {
			flag.PrintDefaults()
			fmt.Println("")
			fmt.Println("Remove source file error:")
			fmt.Println(err)
			os.Exit(1)
			return
		}

		fmt.Println("Locked again!")
		os.Exit(0)
		return
	}
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
