package main

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/cheveuxdelin/keychain/crypt"
	"github.com/cheveuxdelin/keychain/secret"
)

func TestEmptySecret(t *testing.T) {
	var emptySecret secret.Secret = make([]byte, secret.SECRET_SIZE)
	result, _ := secret.CreateSecret([]byte{})
	if !emptySecret.Equals(result) {
		t.Error()
	}
}

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

func TestSuccessfulLogin(t *testing.T) {
	defer os.Remove("test")

	k := createTestKeychain()
	randomSecret := randomSecret()
	k.setPassword(randomSecret)
	k2 := createTestKeychain()
	err := k2.login(randomSecret)
	if err != nil {
		t.Error(err)
	}
}

func TestBadLogin(t *testing.T) {
	defer os.Remove("test")

	k := createTestKeychain()
	k.setPassword(randomSecret())
	k2 := createTestKeychain()
	err := k2.login(randomSecret())
	if err == nil {
		t.Error(err)
	}
}

func TestCredentials(t *testing.T) {

	defer os.Remove("test")
	var testCredentials []credential = []credential{
		{user: "adolfo", password: "adolfo"},
		{user: "cachis", password: "chakis"},
	}

	secret := randomSecret()
	k := createTestKeychain()
	k.setPassword(secret)
	k.credentials = testCredentials
	k.save()

	k2 := createTestKeychain()
	err := k2.login(secret)

	if err != nil {
		t.Error(err)
	}
	if len(testCredentials) != len(k2.credentials) {
		t.Error("error")
	}
	for i := range testCredentials {
		if testCredentials[i].user != k2.credentials[i].user ||
			testCredentials[i].password != k2.credentials[i].password {
			t.Error("error")
		}
	}
}

func randomSecret() (s secret.Secret) {
	s = make([]byte, secret.SECRET_SIZE)
	rand.Read(s)
	return
}

func createTestKeychain() (k Keychain) {
	k = Keychain{
		settings: settings{
			filename:      "test",
			wordDelimiter: ',',
			lineDelimiter: '\n',
		},
	}
	return
}
