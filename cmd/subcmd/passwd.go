package subcmd

import (
	"fmt"

	"coreunit.net/vault/cmd/config"
)

func PasswdOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	fmt.Println("Change password:")

	UnlockOperation(targetFile, appConfig)
	lastUsedPassword = ""
	LockOperation(targetFile, appConfig)

	fmt.Println("Password changed!")
}
