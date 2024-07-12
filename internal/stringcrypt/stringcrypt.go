package stringcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"strings"
)

func AesEncrypt(key string, sourceString string) (string, error) {
	keyBytes := []byte(key)
	sourceBytes := []byte(sourceString)

	keyBytes = IncreaseAesKeyLength(keyBytes, 32)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	bytes := base64.StdEncoding.EncodeToString(sourceBytes)
	cipherSourceBytes := make([]byte, aes.BlockSize+len(bytes))
	iv := cipherSourceBytes[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherSourceBytes[aes.BlockSize:], []byte(bytes))
	return string(cipherSourceBytes), nil
}

func AesDecrypt(key string, encodedString string) (string, error) {
	keyBytes := []byte(key)
	encodedBytes := []byte(encodedString)

	keyBytes = IncreaseAesKeyLength(keyBytes, 32)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	if len(encodedBytes) < aes.BlockSize {
		return "", errors.New("aes cipher encoded string is too short")
	}
	iv := encodedBytes[:aes.BlockSize]
	encodedBytes = encodedBytes[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(encodedBytes, encodedBytes)
	data, err := base64.StdEncoding.DecodeString(string(encodedBytes))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func IncreaseAesKeyLength(key []byte, keySize int) []byte {
	for len(key) < keySize {
		key = IncreaseAesKeyLength(append(key, key...), keySize)
	}

	for len(key) > keySize {
		key = key[1:]
	}

	return key
}

func X509Encrypt(publicKeyPem string, sourceString string) (string, error) {
	publicKeyBlock, _ := pem.Decode([]byte(publicKeyPem))

	if publicKeyBlock == nil {
		return "", errors.New("failed to parse PEM block containing the public key:\n'" + publicKeyPem + "'")
	}

	if publicKeyBlock.Type != "RSA PUBLIC KEY" {
		return "", errors.New("wrong public key type")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return "", errors.New("failed to parse encoded public key:\n> " + err.Error())
	}

	encodedString, err := rsa.EncryptPKCS1v15(
		rand.Reader,
		publicKey,
		[]byte(sourceString),
	)
	if err != nil {
		return "", errors.New("failed to encrypt string:\n> " + err.Error())
	}

	return string(encodedString), nil
}

func X509Decrypt(privateKeyPem string, encodedString string) (string, error) {
	privateKeyPem = strings.TrimSpace(privateKeyPem)

	privateKeyBlock, _ := pem.Decode([]byte(privateKeyPem))

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return "", errors.New("failed to parse encoded private key:\n> " + err.Error())
	}
	decodedString, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(encodedString))
	if err != nil {
		return "", errors.New("failed to decrypt string:\n> " + err.Error())
	}

	return string(decodedString), nil
}
