package utils

import (
	"testing"
)

func TestEmptySecret(t *testing.T) {
	var emptySecret Secret = make([]byte, SECRET_SIZE)
	result, _ := CreateSecret([]byte{})
	if !emptySecret.Equals(result) {
		t.Error()
	}
}
