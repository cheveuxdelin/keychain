package crypt_test

import (
	"testing"

	"github.com/cheveuxdelin/keychain/crypt"
	"github.com/cheveuxdelin/keychain/secret"
)

func TestEncryptAndDecrpyt(t *testing.T) {
	secret_test, err := secret.CreateSecret([]byte("test_secret"))
	data_to_test := []byte("holis")
	if err != nil {
		t.Error(err)
	}
	encrypted, err := crypt.Encrypt(data_to_test, secret_test)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := crypt.Decrypt(encrypted, secret_test)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < len(data_to_test); i++ {
		if data_to_test[i] != decrypted[i] {
			t.Error("el pepe")
		}
	}
}
