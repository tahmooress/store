package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
)

var (
	ErrPasswordIsWrong = errors.New("password is wrong")
	ErrDecrypt         = errors.New("error for decrypting password")
)

func encryptWithAES(password, data []byte) ([]byte, error) {
	salt, err := hash(password)
	if err != nil {
		return nil, fmt.Errorf("encryptWithAES: %s", err)
	}

	cf, err := aes.NewCipher(salt)
	if err != nil {
		return nil, fmt.Errorf("encryptWithAES: %s", err)
	}

	gcm, err := cipher.NewGCM(cf)
	if err != nil {
		return nil, fmt.Errorf("encryptWithAES: %s", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("encryptWithAES: %s", err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decryptAES(cipherText, password []byte) ([]byte, error) {
	salt, err := hash(password)
	if err != nil {
		return nil, fmt.Errorf("decryptAES: %s", err)
	}

	c, err := aes.NewCipher(salt)
	if err != nil {
		return nil, fmt.Errorf("decryptAES: %s", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("decryptAES: %s", err)
	}

	nonceSize := gcm.NonceSize()

	if len(cipherText) < nonceSize {
		return nil, ErrDecrypt
	}

	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrPasswordIsWrong
	}

	return plaintext, nil
}

func hash(data []byte) ([]byte, error) {
	h := sha256.New()

	_, err := h.Write(data)
	if err != nil {
		return nil, fmt.Errorf("hash: %s", err)
	}

	return h.Sum(nil), nil
}
