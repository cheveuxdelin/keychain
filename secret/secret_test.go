package secret_test

import (
	"testing"

	"github.com/cheveuxdelin/keychain/secret"
)

func TestEmptySecret(t *testing.T) {
	var emptySecret secret.Secret = make([]byte, secret.SECRET_SIZE)
	result, _ := secret.CreateSecret([]byte{})
	if !emptySecret.Equals(result) {
		t.Error()
	}
}
