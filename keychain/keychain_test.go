package keychain

import (
	"os"
	"testing"

	"github.com/cheveuxdelin/keychain/utils"
)

func TestCredentials(t *testing.T) {
	defer os.Remove("test")
	var testCredentials []credential = []credential{
		{user: "adolfo", password: "adolfo"},
		{user: "cachis", password: "chakis"},
	}
	secret := utils.RandomSecret()
	k := CreateTestKeychain()
	k.setPassword(secret)
	k.credentials = testCredentials
	k.save()

	k2 := CreateTestKeychain()
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

func TestSuccessfulLogin(t *testing.T) {
	defer os.Remove("test")
	k := CreateTestKeychain()
	randomSecret := utils.RandomSecret()
	k.setPassword(randomSecret)
	k2 := CreateTestKeychain()
	err := k2.login(randomSecret)
	if err != nil {
		t.Error(err)
	}
}

func TestBadLogin(t *testing.T) {
	defer os.Remove("test")
	k := CreateTestKeychain()
	k.setPassword(utils.RandomSecret())
	k2 := CreateTestKeychain()
	err := k2.login(utils.RandomSecret())
	if err == nil {
		t.Error(err)
	}
}
