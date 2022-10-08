package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// http://www.inanzzz.com/index.php/post/f3pe/data-encryption-and-decryption-with-a-secret-key-in-golang
// https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/

// encrypt encrypts plain string with a secret key and returns encrypt string.
func Encrypt(plainData []byte, secret []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(secret)
	if err != nil {
		return []byte{}, err
	}
	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return []byte{}, err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}
	return aead.Seal(nonce, nonce, []byte(plainData), nil), nil
}

// decrypt decrypts encrypt string with a secret key and returns plain string.
func Decrypt(encodedData []byte, secret []byte) ([]byte, error) {
	c, err := aes.NewCipher(secret)
	if err != nil {
		return []byte{}, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return []byte{}, err
	}
	nonceSize := gcm.NonceSize()
	if len(encodedData) < nonceSize {
		return []byte{}, err
	}
	nonce, cipherText := encodedData[:nonceSize], encodedData[nonceSize:]
	plainData, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return []byte{}, err
	}
	return plainData, nil
}
