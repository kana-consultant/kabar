package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

type Decryptor struct {
	key []byte
}

func NewDecryptor(key string) *Decryptor {
	if key == "" {
		key = "01234567890123456789012345678901" // Default key for development
	}
	return &Decryptor{key: []byte(key)}
}

func (d *Decryptor) Decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(d.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("invalid ciphertext")
	}

	nonce, encrypted := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
