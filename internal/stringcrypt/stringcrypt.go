package stringcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh"
)

func deriveKey(passwordBytes []byte, salt []byte, iterations, keySize int) []byte {
	return pbkdf2.Key(passwordBytes, salt, iterations, keySize, sha256.New)
}

func generateHMAC(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func verifyHMAC(key, data, mac []byte) bool {
	expectedMAC := generateHMAC(key, data)
	return hmac.Equal(mac, expectedMAC)
}

func AES256Encrypt(password string, plainPayload string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	keyBytes := deriveKey(
		[]byte(password),
		salt,
		4096,
		32,
	)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	plainBytes := []byte(plainPayload)
	cipherText := make([]byte, aes.BlockSize+len(plainBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainBytes)

	mac := generateHMAC(keyBytes, cipherText)
	finalCipherText := append(salt, append(cipherText, mac...)...)

	return base64.StdEncoding.EncodeToString(finalCipherText), nil
}

func AES256Decrypt(password string, cipherPayload string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherPayload)
	if err != nil {
		return "", err
	}

	if len(cipherText) < 16+aes.BlockSize+sha256.Size {
		return "", errors.New("cipher text too short")
	}

	salt := cipherText[:16]
	cipherText = cipherText[16:]

	keyBytes := deriveKey([]byte(password), salt, 4096, 32)

	// Separate the HMAC from the ciphertext
	hmacStart := len(cipherText) - sha256.Size
	mac := cipherText[hmacStart:]
	cipherText = cipherText[:hmacStart]

	if !verifyHMAC(keyBytes, cipherText, mac) {
		return "", errors.New("decryption failed: invalid key or corrupted data")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipher text too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func X509Encrypt(publicKey *rsa.PublicKey, plainPayload string) (string, error) {
	if len(plainPayload) == 0 {
		return "", errors.New("empty plain payload")
	} else if publicKey == nil {
		return "", errors.New("nil public key")
	}

	encodedString, err := rsa.EncryptPKCS1v15(
		rand.Reader,
		publicKey,
		[]byte(plainPayload),
	)
	if err != nil {
		return "", errors.New("failed to encrypt string:\n> " + err.Error())
	}

	return string(encodedString), nil
}

func X509Decrypt(privateKey *rsa.PrivateKey, cipherPayload string) (string, error) {
	if len(cipherPayload) == 0 {
		return "", errors.New("empty cipher payload")
	} else if privateKey == nil {
		return "", errors.New("nil private key")
	}

	decodedString, err := rsa.DecryptPKCS1v15(
		rand.Reader,
		privateKey,
		[]byte(cipherPayload),
	)
	if err != nil {
		return "", errors.New("failed to decrypt string:\n> " + err.Error())
	}

	return string(decodedString), nil
}

func LoadRsaPublicKey(path string) (*rsa.PublicKey, error) {
	filePayload, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("error while read public key file:\n> " + err.Error())
	}

	fileContent := strings.TrimSpace(string(filePayload))
	filePayload = []byte(fileContent)

	var publicKey *rsa.PublicKey

	if strings.HasPrefix(fileContent, "ssh-rsa") {
		sshPublicKey, _, _, _, err := ssh.ParseAuthorizedKey(filePayload)
		if err != nil {
			return nil, errors.New("error parsing authorized public key:\n> " + err.Error())
		}

		parsedCryptoKey, ok := sshPublicKey.(ssh.CryptoPublicKey)
		if !ok {
			return nil, errors.New("unsupported parsed ssh public authorized key type, need to be openssh authorized key or pem encoded rsa public key")
		}
		pubCrypto := parsedCryptoKey.CryptoPublicKey()
		publicKey, ok = pubCrypto.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("unsupported crypto public key type, need to be openssh authorized key or pem encoded rsa public key")
		}
	} else if strings.HasPrefix(fileContent, "-----BEGIN PUBLIC KEY-----") ||
		strings.HasPrefix(fileContent, "-----BEGIN RSA PUBLIC KEY-----") {

		pemBlock, _ := pem.Decode(filePayload)
		publicKey, err = x509.ParsePKCS1PublicKey(pemBlock.Bytes)
		if err != nil {
			return nil, errors.New("error parsing rsa public key:\n> " + err.Error())
		}
	} else {
		return nil, errors.New("unsupported public key format, need to be openssh authorized key or pem encoded rsa public key")
	}

	return publicKey, nil
}

func LoadRsaPrivateKey(path string) (*rsa.PrivateKey, error) {
	filePayload, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("error while read private key file:\n> " + err.Error())
	}

	fileContent := strings.TrimSpace(string(filePayload))
	filePayload = []byte(fileContent)

	var privateKey *rsa.PrivateKey

	if strings.HasPrefix(fileContent, "-----BEGIN PRIVATE KEY-----") ||
		strings.HasPrefix(fileContent, "-----BEGIN RSA PRIVATE KEY-----") {
		pemBlock, _ := pem.Decode(filePayload)
		if pemBlock == nil {
			return nil, errors.New("failed to decode PEM block containing rsa private key")
		}
		genericPrivateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		if err != nil {
			return nil, errors.New("error parsing rsa private key:\n> " + err.Error())
		}
		privateKey = genericPrivateKey
	} else if strings.HasPrefix(fileContent, "-----BEGIN OPENSSH PRIVATE KEY-----") {
		key, err := ssh.ParseRawPrivateKey(filePayload)
		if err != nil {
			return nil, errors.New("error parsing openssh private key:\n> " + err.Error())
		}

		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("unsupported openssh private key type, need to be openssh or pem encoded rsa private key")
		}
	} else {
		return nil, errors.New("unsupported private key format, need to be openssh or pem encoded rsa private key")
	}

	return privateKey, nil
}
