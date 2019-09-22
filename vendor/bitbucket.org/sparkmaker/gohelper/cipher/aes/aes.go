package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

type AES struct {
	hash []byte
}

func New(passphrase string) *AES {
	return &AES{hash: []byte(createHash(passphrase))}
}

func (c *AES) Encrypt(data string) string {
	block, _ := aes.NewCipher(c.hash)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%x", gcm.Seal(nonce, nonce, []byte(data), nil))
}

func (c *AES) Decrypt(encrypted string) (string, error) {
	if len(strings.TrimSpace(encrypted)) == 0 {
		return "", errors.New("cipher text is empty")
	}

	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.hash)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
