package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
)

// Encrypt is using for encrypting data
func Encrypt(password, data string) (string, error) {
	key := sha256.Sum256([]byte(password)) // ключ шифрования

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	return hex.EncodeToString(aesgcm.Seal(nil, nonce, []byte(data), nil)), nil

}

// Decrypt is using for decrypting data
func Decrypt(password string, encryptedData string) (string, error) {
	key := sha256.Sum256([]byte(password)) // ключ шифрования

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	data, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	decrypted, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
