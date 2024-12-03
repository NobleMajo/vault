package main

import (
	"fmt"
	"os"
	"strings"

	"coreunit.net/vault/internal/config"
	"coreunit.net/vault/internal/subcmd"
	"coreunit.net/vault/lib/stringfs"
	 "github.com/joho/godotenv"
)

var DisplayName string = "Unset"
var ShortName string = "unset"
var Version string = "?.?.?"
var Commit string = "???????"

func main() {
	err := godotenv.Load()
	if err == nil {
		fmt.Println("Environment variables from .env loaded")
	}
	
	appConfig := config.ParseConfig(DisplayName, ShortName, Version, Commit)

	stringfs.ParsePath(&appConfig.PublicKeyPath)
	stringfs.ParsePath(&appConfig.PrivateKeyPath)
	targetFile := targetFile(appConfig)

	if appConfig.SubCommand == "lock" {
		subcmd.LockOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "init" {
		subcmd.InitOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "print" {
		subcmd.PrintOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "unlock" {
		subcmd.UnlockOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "temp" {
		subcmd.TempOperation(
			targetFile,
			appConfig,
		)
	} else if appConfig.SubCommand == "passwd" {
		subcmd.PasswdOperation(
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
