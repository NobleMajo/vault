package subcmd

import (
	"fmt"
	"strconv"
	"time"

	"coreunit.net/vault/cmd/config"
)

var lastUsedPassword string

func TempOperation(
	targetFile string,
	appConfig *config.AppConfig,
) {
	fmt.Println(
		"Temporary unlock vault for " +
			strconv.Itoa(appConfig.TempDecodeSeconds) +
			" seconds...",
	)

	UnlockOperation(targetFile, appConfig)

	fmt.Println("Wait for " + strconv.Itoa(appConfig.TempDecodeSeconds) + " seconds...")
	time.Sleep(time.Duration(appConfig.TempDecodeSeconds) * time.Second)
	fmt.Println("Lock vault now...")

	LockOperation(targetFile, appConfig)

	fmt.Println("Locked again!")
}
