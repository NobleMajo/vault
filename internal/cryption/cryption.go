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

func pubicKeyMaxMessageLength(pubicKey *rsa.PublicKey) int {
	return pubicKey.Size() - 11
}

func privateKeyMaxMessageLength(privateKey *rsa.PrivateKey) int {
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

func RandomByteArray(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	return randomBytes, nil
}

func AES256Encrypt(key []byte, plainPayload []byte) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	keyBytes := deriveKey(
		key,
		salt,
		4096,
		32,
	)

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
	if len(cipherPayload) < 16+aes.BlockSize+sha256.Size {
		return nil, errors.New("cipher text too short")
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

func X509AES256Encrypt(publicKey *rsa.PublicKey, plainPayload []byte) ([]byte, error) {
	if len(plainPayload) == 0 {
		return nil, errors.New("empty plain payload")
	} else if publicKey == nil {
		return nil, errors.New("nil public key")
	}

	randomKey, err := RandomByteArray(pubicKeyMaxMessageLength(publicKey))
	if err != nil {
		return nil, err
	}

	result, err := AES256Encrypt(randomKey, plainPayload)
	if err != nil {
		return nil, err
	}

	encryptedKey, err := X509EncryptPart(publicKey, randomKey)
	if err != nil {
		return nil, err
	}

	return append(encryptedKey, result...), nil
}

func X509AES256Decrypt(privateKey *rsa.PrivateKey, cipherPayload []byte) ([]byte, error) {
	if len(cipherPayload) == 0 {
		return nil, errors.New("empty cipher payload")
	} else if privateKey == nil {
		return nil, errors.New("nil private key")
	}

	keySize := privateKey.Size()
	encryptedKey := cipherPayload[:keySize]

	plainKey, err := X509DecryptPart(privateKey, encryptedKey)
	if err != nil {
		return nil, err
	}

	result, err := AES256Decrypt(plainKey, cipherPayload[keySize:])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func X509EncryptPart(publicKey *rsa.PublicKey, plainPayload []byte) ([]byte, error) {
	if len(plainPayload) == 0 {
		return nil, errors.New("empty plain payload")
	} else if publicKey == nil {
		return nil, errors.New("nil public key")
	}

	encodedString, err := rsa.EncryptPKCS1v15(
		rand.Reader,
		publicKey,
		plainPayload,
	)
	if err != nil {
		return nil, errors.New("failed to encrypt string:\n> " + err.Error())
	}

	return encodedString, nil
}

func X509DecryptPart(privateKey *rsa.PrivateKey, cipherPayload []byte) ([]byte, error) {
	if len(cipherPayload) == 0 {
		return nil, errors.New("empty cipher payload")
	} else if privateKey == nil {
		return nil, errors.New("nil private key")
	}

	decodedString, err := rsa.DecryptPKCS1v15(
		nil,
		privateKey,
		cipherPayload,
	)
	if err != nil {
		return nil, errors.New("failed to decrypt string:\n> " + err.Error())
	}

	return decodedString, nil
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
