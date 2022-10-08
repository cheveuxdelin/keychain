package utils

import (
	"errors"
)

type Secret []byte

const SECRET_SIZE int = 32

var errorBadSecretSize error = errors.New("secret must be between 1-32 ASCII characters")

func CreateSecret(b []byte) (s Secret, err error) {
	if len(b)-1 > SECRET_SIZE || len(b)-1 == 0 {
		CheckError(err)
	}
	s = make(Secret, SECRET_SIZE)
	for i := 0; i < len(b)-1; i++ {
		s[i] = b[i]
	}
	return
}

func (s Secret) Equals(s2 Secret) bool {
	if len(s) != len(s2) {
		return false
	}
	for i := 0; i < SECRET_SIZE; i++ {
		if s[i] != s2[i] {
			return false
		}
	}
	return true
}
