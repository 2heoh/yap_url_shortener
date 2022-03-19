package services

import (
	"crypto/aes"
	"crypto/cipher"
)

type crypto struct {
	aesblock cipher.Block
	aesgcm   cipher.AEAD
	nonce    []byte
}

func (c *crypto) Encrypt(id []byte) []byte {
	return c.aesgcm.Seal(nil, c.nonce, id, nil)
}
func (c *crypto) Decrypt(src []byte) ([]byte, error) {
	return c.aesgcm.Open(nil, c.nonce, src, nil) // расшифровываем
}

func generateNonce(size int) []byte {
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b[i] = byte(i)
	}
	return b
}

func NewCrypto() (*crypto, error) {
	secretkey := []byte("1234567890abcdef1234567890abcdef")
	aesblock, err := aes.NewCipher(secretkey)
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
