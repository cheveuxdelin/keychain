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
	var c map[string]string = map[string]string{
		"el pepe":    "ete sech",
		"la chona":   "se mueve",
		"y la gente": "le grita",
	}

	secret := randomSecret()
	k := createTestKeychain()
	k.setPassword(secret)
	k.credentials = c
	k.save()

	k2 := createTestKeychain()
	err := k2.login(secret)

	if err != nil {
		t.Error(err)
	}

	if len(k2.credentials) != len(c) {
		t.Error("popo")
	}
	for key, value := range c {
		value2, ok := k2.credentials[key]
		if !ok {
			t.Error("popo")
		}
		if value != value2 {
			t.Error("pop2")
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
