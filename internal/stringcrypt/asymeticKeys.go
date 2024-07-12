package stringcrypt

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"coreunit.net/vault/internal/stringfs"
	"golang.org/x/crypto/ssh"
)

func ParseAnyToPemPrivateKey(anyPrivateKey string, getPassphrase func() (string, error)) (string, error) {
	anyPrivateKey = strings.TrimSpace(anyPrivateKey)
	keyData := []byte(anyPrivateKey)

	var der []byte

	if strings.HasPrefix(anyPrivateKey, "-----BEGIN PRIVATE KEY-----") {
		privateKey, err := x509.ParsePKCS1PrivateKey(keyData)
		if err != nil {
			return "", fmt.Errorf("failed to parse PKCS#1 private key:\n> %v", err)
		}

		der = x509.MarshalPKCS1PrivateKey(privateKey)
	} else if strings.HasPrefix(anyPrivateKey, "-----BEGIN EC PRIVATE KEY-----") {
		ecKey, err := x509.ParseECPrivateKey(keyData)
		if err != nil {
			return "", fmt.Errorf("failed to parse EC private key:\n> %v", err)
		}
		der, err = x509.MarshalECPrivateKey(ecKey)
		if err != nil {
			return "", fmt.Errorf("failed to marshal EC private key:\n> %v", err)
		}
	} else if strings.HasPrefix(anyPrivateKey, "-----BEGIN ENCRYPTED PRIVATE KEY-----") {
		passphrase, err := getPassphrase()
		if err != nil {
			return "", fmt.Errorf("failed to get passphrase for an encrypted private key:\n> %v", err)
		}

		block, rest := pem.Decode(keyData)
		if len(rest) > 0 {
			return "", fmt.Errorf("extra data included in key, expected no more data")
		}
		der, err = x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return "", fmt.Errorf("failed to decrypt private key:\n> %v", err)
		}
	} else if strings.HasPrefix(anyPrivateKey, "-----BEGIN OPENSSH PRIVATE KEY-----") {
		key, err := ssh.ParseRawPrivateKey(keyData)
		if err != nil {
			return "", fmt.Errorf("failed to parse OpenSSH private key:\n> %v", err)
		}

		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("the provided key is not a RSA private key")
		}

		der = x509.MarshalPKCS1PrivateKey(rsaKey)
	} else {
		return "", errors.New("invalid private key format")
	}

	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	}

	var pemBuffer bytes.Buffer
	if err := pem.Encode(&pemBuffer, pemBlock); err != nil {
		return "", fmt.Errorf("failed to encode PEM block:\n> %v", err)
	}

	return pemBuffer.String(), nil
}

func ConvertOpensshPublicKeyToPemRsaPublicKey(publicKey string) (string, error) {
	var parsedPubKey ssh.PublicKey
	var err error

	if strings.HasPrefix(publicKey, "ssh-") {
		parsedPubKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte(publicKey))
	} else {
		parsedPubKey, err = ssh.ParsePublicKey([]byte(publicKey))
	}
	if err != nil {
		return "", errors.New("invalid OpenSSH public key format:\n> " + err.Error())
	}

	cryptoKey, ok := parsedPubKey.(ssh.CryptoPublicKey)
	if !ok {
		return "", errors.New("invalid OpenSSH public key format")
	}

	rsaKey := cryptoKey.CryptoPublicKey().(*rsa.PublicKey)
	derBytes := x509.MarshalPKCS1PublicKey(rsaKey)

	pemBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derBytes,
	}

	pemKey := pem.EncodeToMemory(&pemBlock)
	return string(pemKey), nil
}

func LoadPublicKey(keyDir string, publicKeyNames []string) (string, error) {
	var err2 error = errors.New("no public key found in " + keyDir)

	for _, publicKeyName := range publicKeyNames {
		publicKeyPath := keyDir + "/" + publicKeyName
		_, isFile := stringfs.IsFile(publicKeyPath)
		if isFile {
			content, err := stringfs.ReadFile(publicKeyPath)

			if err != nil {
				err2 = err
				continue
			}

			content = strings.TrimSpace(content)

			if !strings.HasPrefix(content, "-----BEGIN") {
				content, err = ConvertOpensshPublicKeyToPemRsaPublicKey(content)

				if err != nil {
					return "", errors.New("failed to convert OpenSSH public key to PEM format:\n> " + err.Error())
				}
			}

			publicKeyBlock, _ := pem.Decode([]byte(content))

			if publicKeyBlock == nil {
				return "", errors.New("failed to parse PEM block containing the public key:\n'" + content + "'")
			}

			if publicKeyBlock.Type != "RSA PUBLIC KEY" {
				return "", errors.New("wrong X509 public key type: " + publicKeyBlock.Type)
			}

			return content, nil
		}
	}

	return "", err2
}

func LoadPrivateKey(keyDir string, privateKeyNames []string, getPassphrase func() (string, error)) (string, error) {
	var err2 error

	for _, privateKeyName := range privateKeyNames {
		privateKeyPath := keyDir + "/" + privateKeyName
		_, isFile := stringfs.IsFile(privateKeyPath)
		if isFile {
			content, err := stringfs.ReadFile(privateKeyPath)
			if err != nil {
				err2 = err
				continue
			}

			content, err = ParseAnyToPemPrivateKey(content, getPassphrase)

			if err == nil {
				return content, nil
			}

			err2 = fmt.Errorf("failed to parse private key:\n> %v", err)
		}
	}

	return "", err2
}
