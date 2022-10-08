package utils

import (
	"bufio"
	"crypto/rand"
	"log"
	"os"
	"strings"

	"github.com/cheveuxdelin/keychain/secret"
	"golang.org/x/term"
)

func IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func ReadSafeBytes() (b []byte) {
	b, err := term.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func GetEnteredPassword() secret.Secret {
	b := ReadSafeBytes()
	secret, err := secret.CreateSecret(b)
	if err != nil {
		log.Fatal(err)
	}
	return secret
}

func ReadBytes() (b []byte) {
	reader := bufio.NewReader(os.Stdin)
	b, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	return
}

func ReadString() (s string) {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(s)
}

func RandomSecret() (s secret.Secret) {
	s = make([]byte, secret.SECRET_SIZE)
	rand.Read(s)
	return
}
