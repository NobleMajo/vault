package subcmd

import (
	"crypto/rsa"
	"fmt"
	"os"

	"coreunit.net/vault/internal/config"
	"coreunit.net/vault/lib/cryption"
	"coreunit.net/vault/lib/userin"
)

var err error
var lastUsedPrivateKey *rsa.PrivateKey
var lastUsedPublicKey *rsa.PublicKey
var lastUsedPassword string

func loadDecryptionData(appConfig *config.AppConfig) {
	if !appConfig.DisableRSA {
		lastUsedPrivateKey, err = cryption.LoadRsaPrivateKey(appConfig.PrivateKeyPath)

		if err != nil {
			exitError("Load private key error:\n> " + err.Error())
			return
		}
	}

	if !appConfig.DisableAES256 && len(lastUsedPassword) == 0 {
		lastUsedPassword, err = userin.PromptPassword()

		if err != nil {
			exitError("Prompt new password error:\n> " + err.Error())
			return
		}
	}
}

func loadEncryptionData(appConfig *config.AppConfig) {
	if !appConfig.DisableRSA {
		lastUsedPublicKey, err = cryption.LoadRsaPublicKey(appConfig.PublicKeyPath)

		if err != nil {
			exitError("Load public key error:\n> " + err.Error())
			return
		}
	}

	if !appConfig.DisableAES256 && len(lastUsedPassword) == 0 {
		lastUsedPassword, err = userin.PromptNewPassword()

		if err != nil {
			exitError("Prompt new password error:\n> " + err.Error())
			return
		}
	}
}

func exitError(message string) {
	fmt.Fprintln(
		os.Stderr,
		message,
	)
	os.Exit(1)
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
