package main

import (
	"fmt"
	"os"
	"strings"

	"coreunit.net/vault/cmd/config"
	"coreunit.net/vault/cmd/subcmd"
	"coreunit.net/vault/internal/stringfs"
)

var Version string = "?.?.?"
var Commit string = "???????"

func main() {
	if len(Commit) > 7 {
		Commit = Commit[0:7]
	}

	appConfig := config.ParseConfig(Version, Commit)

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
