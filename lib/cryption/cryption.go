package cryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh"
)

func pubicKeyMaxMessageLength(publicKey *rsa.PublicKey) int {
	if publicKey == nil {
		return 0
	}
	return publicKey.Size() - 11
}

func privateKeyMaxMessageLength(privateKey *rsa.PrivateKey) int {
	if privateKey == nil {
		return 0
	}
	return privateKey.Size() - 11
}

func splitByteSliceIntoSize(data []byte, size int) [][]byte {
	if size <= 0 {
		return nil
	}

	var result [][]byte
	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		result = append(result, data[i:end])
	}
	return result
}

func deriveKey(passwordBytes []byte, salt []byte, iterations, keySize int) []byte {
	if len(passwordBytes) == 0 {
		panic("password cannot be empty")
	}
	if len(salt) == 0 {
		panic("salt cannot be empty")
	}
	if iterations <= 0 || keySize <= 0 {
		panic("invalid iterations or key size")
	}
	return pbkdf2.Key(passwordBytes, salt, iterations, keySize, sha256.New)
}

func generateHMAC(key, data []byte) []byte {
	if len(key) == 0 {
		panic("HMAC key cannot be empty")
	}
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func verifyHMAC(key, data, mac []byte) bool {
	if len(key) == 0 || len(mac) == 0 {
		return false
	}
	expectedMAC := generateHMAC(key, data)
	return hmac.Equal(mac, expectedMAC)
}

func RandomByteArray(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("length must be positive")
	}
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}

func AES256Encrypt(key []byte, plainPayload []byte) ([]byte, error) {
	if len(key) == 0 || len(plainPayload) == 0 {
		return nil, errors.New("key and plainPayload cannot be empty")
	}

	salt, err := RandomByteArray(16)
	if err != nil {
		return nil, err
	}

	keyBytes := deriveKey(key, salt, 4096, 32)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	cipherText := make([]byte, aes.BlockSize+len(plainPayload))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainPayload)

	mac := generateHMAC(keyBytes, cipherText)
	finalCipherText := append(salt, append(cipherText, mac...)...)

	return finalCipherText, nil
}

func AES256Decrypt(key []byte, cipherPayload []byte) ([]byte, error) {
	if len(key) == 0 || len(cipherPayload) < 16+aes.BlockSize+sha256.Size {
		return nil, errors.New("invalid input data for decryption")
	}

	salt := cipherPayload[:16]
	cipherPayload = cipherPayload[16:]

	keyBytes := deriveKey(key, salt, 4096, 32)

	hmacStart := len(cipherPayload) - sha256.Size
	mac := cipherPayload[hmacStart:]
	cipherPayload = cipherPayload[:hmacStart]

	if !verifyHMAC(keyBytes, cipherPayload, mac) {
		return nil, errors.New("decryption failed: invalid key or corrupted data")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	if len(cipherPayload) < aes.BlockSize {
		return nil, errors.New("cipher text too short")
	}
	iv := cipherPayload[:aes.BlockSize]
	cipherPayload = cipherPayload[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherPayload, cipherPayload)

	return cipherPayload, nil
}

func X509EncryptPart(publicKey *rsa.PublicKey, plainPayload []byte) ([]byte, error) {
	if len(plainPayload) == 0 || publicKey == nil {
		return nil, errors.New("publicKey and plainPayload must be non-nil and non-empty")
	}

	if len(plainPayload) > pubicKeyMaxMessageLength(publicKey) {
		return nil, errors.New("plainPayload too large for encryption")
	}

	encodedString, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainPayload)
	if err != nil {
		return nil, errors.New("failed to encrypt string:\n> " + err.Error())
	}

	return encodedString, nil
}

func X509DecryptPart(privateKey *rsa.PrivateKey, cipherPayload []byte) ([]byte, error) {
	if len(cipherPayload) == 0 || privateKey == nil {
		return nil, errors.New("privateKey and cipherPayload must be non-nil and non-empty")
	}

	decodedString, err := rsa.DecryptPKCS1v15(nil, privateKey, cipherPayload)
	if err != nil {
		return nil, errors.New("failed to decrypt string:\n> " + err.Error())
	}

	return decodedString, nil
}

// Key-loading functions omitted for brevity. Add similar validations for file existence, format correctness, and content handling.
