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

func X509PubicKeyMaxEncryptPayloadLength(pubicKey *rsa.PublicKey) int {
	return pubicKey.Size() - 11
}

func X509PrivateKeyMaxEncryptPayloadLength(privateKey *rsa.PrivateKey) int {
	return privateKey.Size() - 11
}

func SplitByteSliceIntoSize(data []byte, size int) [][]byte {
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
	if key == nil ||
	len(key) == 0 ||
	data == nil ||
	len(data) == 0 ||
	mac == nil ||
	len(mac) == 0 {
		return false
	}
	expectedMAC := generateHMAC(key, data)
	return hmac.Equal(mac, expectedMAC)
}

// RandomByteArray returns a byte slice of length `length` that is randomly generated.
//
func RandomByteArray(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	return randomBytes, nil
}

// AES256Encrypt encrypts the given plainPayload with the given key. It returns the
// encrypted text or an error. The encrypted text is a concatenation of the
// random salt used in the derivation of the encryption key, the encrypted
// payload, and the MAC (Message Authentication Code) of the encrypted payload.
// The salt is a random byte slice of length 16, the encrypted payload is a
// byte slice of length BlockSize + len(plainPayload) where BlockSize is the
// block size of the AES cipher, and the MAC is a byte slice of length
// sha256.Size (32 bytes).
func AES256Encrypt(key []byte, plainPayload []byte) ([]byte, error) {
	if key == nil {
		return nil, errors.New("nil key")
	}else if len(key) == 0 {
		return nil, errors.New("empty key")
	} else if plainPayload == nil {
		return nil, errors.New("nil plain payload")
	} else if len(plainPayload) == 0 {
		return nil, errors.New("empty plain payload")
	}

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

// AES256Decrypt takes a key and cipher text and returns the decrypted plain text.
//
// The cipher text is expected to contain a 16 byte salt, followed by the AES-256
// encrypted plain text, followed by the 256 bit HMAC of the plain text.
//
// The key is expected to be a 256 bit key, it is used to derive a 256 bit key
// using the salt and the PBKDF2 key derivation function with a cost of 4096.
//
// The HMAC is verified before the plain text is decrypted, if the HMAC is invalid
// an error is returned.
//
// The cipher text is decrypted using the derived key and the AES-256 cipher in
// Cipher Feedback (CFB) mode.
//
// The decrypted plain text is returned as a byte slice, or an error is returned
// if any of the above steps fail.
func AES256Decrypt(key []byte, cipherPayload []byte) ([]byte, error) {
	if key == nil {
		return nil, errors.New("nil key")
	}else if len(key) == 0 {
		return nil, errors.New("empty key")
	} else if cipherPayload == nil {
		return nil, errors.New("nil cipher text")
	} else if len(cipherPayload) == 0 {
		return nil, errors.New("empty cipher text")
	}

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

// Requires a rsa public key and a payload (that should be encrypted) as parameters.
// Generates a random byte array with the maximum encryption length technically permitted by the public key.
// Then uses this random byte array to encrypt the payload via aes encryption and returns the encrypted playload (cipher payload) as well as the X509 encrypted random byte array.
// The cipher payload is a concatenation of a X509 encrypted random byte array and the aes encrypted payload.
// 
// Even if a symmetric aes encryption is used internally, the actual procedure must be regarded as asymmetric because only the private key can decrypt the random byte array, which is the only one that can decrypt the playload via aes.
func X509AES256Encrypt(publicKey *rsa.PublicKey, plainPayload []byte) ([]byte, error) {
	if len(plainPayload) == 0 {
		return nil, errors.New("empty plain payload")
	} else if publicKey == nil {
		return nil, errors.New("nil public key")
	}

	randomKey, err := RandomByteArray(X509PubicKeyMaxEncryptPayloadLength(publicKey))
	if err != nil {
		return nil, err
	}

	result, err := AES256Encrypt(randomKey, plainPayload)
	if err != nil {
		return nil, err
	}

	encryptedKey, err := X509ChunkEncrypt(publicKey, randomKey)
	if err != nil {
		return nil, err
	}

	return append(encryptedKey, result...), nil
}

// Requires a rsa private key and a cipher payload (created by X509AES256Encrypt) as parameters.
// First splits the cipher payload into the X509 encrypted random byte array and the aes encrypted cipher payload.
// Then uses the X509 private key to decrypt the X509 encrypted random byte array.
// Then uses the decrypted random byte array to aes decrypt the rest of the cipher payload.
func X509AES256Decrypt(privateKey *rsa.PrivateKey, cipherPayload []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("nil private key")
	} else if cipherPayload == nil {
		return nil, errors.New("nil cipher payload")
	}else 	if len(cipherPayload) == 0 {
		return nil, errors.New("empty cipher payload")
	}  

	keySize := privateKey.Size()
	encryptedKey := cipherPayload[:keySize]

	plainKey, err := X509ChunkDecrypt(privateKey, encryptedKey)
	if err != nil {
		return nil, err
	}

	result, err := AES256Decrypt(plainKey, cipherPayload[keySize:])
	if err != nil {
		return nil, err
	}

	return result, nil
}

// X509ChunkEncrypt encrypts a plain text using a public key. The plain text is expected to be not longer than the maximum payload length depending on the public key size.
//
// X509ChunkEncrypt and X509ChunkDecrypt are called "chunk" encrypt / decrypt because they have a maximum payload length depending on the public key size.
//
// The function will return an error if the public key is nil, the plain payload is nil or empty, or if an error occurs during the encryption process.
func X509ChunkEncrypt(publicKey *rsa.PublicKey, plainPayload []byte) ([]byte, error) {
	if len(plainPayload) == 0 {
		return nil, errors.New("empty plain payload")
	} else if publicKey == nil {
		return nil, errors.New("nil public key")
	}

	if len(plainPayload) > X509PubicKeyMaxEncryptPayloadLength(publicKey) {
		return nil, errors.New("plain payload too long to X509 encrypt")
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

// X509ChunkDecrypt decrypts a ciphertext using a private key. The ciphertext is expected to have been encrypted using the corresponding public key.
//
// X509ChunkEncrypt and X509ChunkDecrypt are called "chunk" encrypt / decrypt because they have a maximum payload length depending on the public key size.
//
// The function will return an error if the private key is nil, the cipher payload is nil or empty, or if an error occurs during the decryption process.
func X509ChunkDecrypt(privateKey *rsa.PrivateKey, cipherPayload []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("nil private key")
	} else if cipherPayload == nil {
		return nil, errors.New("nil cipher payload")
	} else 	if len(cipherPayload) == 0 {
		return nil, errors.New("empty cipher payload")
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
	if len(path) == 0 {
		return nil, errors.New("empty public key path")
	}

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
	if len(path) == 0 {
		return nil, errors.New("empty public key path")
	}

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
