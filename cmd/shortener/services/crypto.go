package services

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"math/rand"
)

type crypto struct {
	aesblock cipher.Block
	aesgcm   cipher.AEAD
	nonce    []byte
}

const (
	userIdLength = 16
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	secretKey    = "1234567890abcdef1234567890abcdef"
)

func NewCrypto() (*crypto, error) {
	aesblock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	var nonce = generateNonce(aesgcm.NonceSize())

	return &crypto{
		aesblock: aesblock,
		aesgcm:   aesgcm,
		nonce:    nonce,
	}, nil
}

func (c *crypto) encrypt(id []byte) []byte {

	return c.aesgcm.Seal(nil, c.nonce, id, nil)
}

func (c *crypto) decrypt(src []byte) ([]byte, error) {

	return c.aesgcm.Open(nil, c.nonce, src, nil) // расшифровываем
}

func generateNonce(size int) []byte {
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b[i] = byte(i)
	}
	return b
}

func (c *crypto) generateUserID() string {
	return c.generateRandomString(userIdLength)
}

func (c *crypto) generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (c *crypto) GetEncodedSessionValue() string {
	UserID := c.generateUserID()
	seal := c.encrypt([]byte(UserID))

	return hex.EncodeToString(seal)
}

func (c *crypto) GetDecodedUserID(rawCookieValue string) (string, error) {
	decodedCookieValue, err := hex.DecodeString(rawCookieValue)
	if err != nil {
		return "", err
	}

	decryptedUserID, err := c.decrypt(decodedCookieValue)
	if err != nil {
		return "", err
	}

	return string(decryptedUserID), nil
}
