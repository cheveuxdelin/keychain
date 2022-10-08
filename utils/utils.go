package utils

import (
	"bufio"
	"crypto/rand"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/eiannone/keyboard"
	"golang.org/x/term"
)

func IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func ReadSafeBytes() (b []byte) {
	b, err := term.ReadPassword(int(syscall.Stdin))
	CheckError(err)
	return
}

func GetEnteredPassword() Secret {
	b := ReadSafeBytes()
	secret, err := CreateSecret(b)
	CheckError(err)
	return secret
}

func ReadBytes() (b []byte) {
	reader := bufio.NewReader(os.Stdin)
	b, err := reader.ReadBytes('\n')
	CheckError(err)
	return
}

func ReadString() (s string) {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	CheckError(err)
	return strings.TrimSpace(s)
}

func RandomSecret() (s Secret) {
	s = make([]byte, SECRET_SIZE)
	_, err := rand.Read(s)
	CheckError(err)
	return
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func CheckError(err error) {
	if err != nil {
		keyboard.Close()
		log.Fatal(err)
	}
}
