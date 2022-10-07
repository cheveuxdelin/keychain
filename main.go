package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cheveuxdelin/keychain/crypt"
	"github.com/cheveuxdelin/keychain/secret"
	"github.com/eiannone/keyboard"
	"github.com/gookit/color"
	"github.com/k0kubun/go-ansi"
)

type settings struct {
	filename      string
	wordDelimiter byte
	lineDelimiter byte
}

type credential struct {
	user     string
	password string
}

type Keychain struct {
	credentials []credential
	secret      secret.Secret
	settings    settings
}

func (k *Keychain) credentialsToBytes() []byte {
	var b bytes.Buffer
	for i := range k.credentials {
		b.WriteString(k.credentials[i].user)
		b.WriteByte(k.settings.wordDelimiter)
		b.WriteString(k.credentials[i].password)
		b.WriteByte(k.settings.lineDelimiter)
	}
	return b.Bytes()
}

func (k *Keychain) save() {
	encrypted, err := crypt.Encrypt(k.credentialsToBytes(), k.secret)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(k.settings.filename, encrypted, 0777)
}

func (k *Keychain) load() (err error) {
	bytes, err := os.ReadFile(k.settings.filename)
	if err != nil {
		return
	}
	decrypted, err := crypt.Decrypt(bytes, k.secret)
	if err != nil {
		return
	}

	k.credentials = []credential{}
	var user strings.Builder
	var password strings.Builder

	for i := 0; i < len(decrypted); {
		for ; decrypted[i] != k.settings.wordDelimiter; i++ {
			user.WriteByte(decrypted[i])
		}
		i++
		for ; decrypted[i] != k.settings.lineDelimiter; i++ {
			password.WriteByte(decrypted[i])
		}
		i++
		k.credentials = append(
			k.credentials,
			credential{
				user:     user.String(),
				password: password.String(),
			})
		user.Reset()
		password.Reset()
	}
	return nil
}

func (k *Keychain) Add(user string, password string) {
	k.credentials = append(k.credentials, credential{user: user, password: password})
	k.save()
}

func (k *Keychain) setPassword(secret secret.Secret) {
	/*
		fmt.Print("Insert new secret 🔑 (Must be between 1-32 ASCII characters): ")
		reader := bufio.NewReader(os.Stdin)
		b, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
	*/
	k.secret = secret
	k.save()
}

func (k *Keychain) login(secret secret.Secret) (err error) {
	/*
			fmt.Print("Insert 🔑: ")
				reader := bufio.NewReader(os.Stdin)
				b, err := reader.ReadBytes('\n')
				if err != nil {
					log.Fatal(err)
				}
				secret, err := secret.CreateSecret(enteredSecret)
		if err != nil {
			return err
		}
	*/
	k.secret = secret
	err = k.load()
	if err != nil {
		return
	}
	return nil
}

func CreateKeychain() (k Keychain) {
	k = Keychain{
		settings: settings{
			filename:      ".credentials",
			wordDelimiter: ',',
			lineDelimiter: '\n',
		},
	}
	return
}

func readLine() (b []byte) {
	reader := bufio.NewReader(os.Stdin)
	b, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (k *Keychain) userAlreadyCreated() bool {
	_, err := os.Stat(k.settings.filename)
	return errors.Is(err, os.ErrNotExist)
}

func (c *credential) print(indexNumber int) string {
	return fmt.Sprintf("[%d] %-20s | %-20s", indexNumber, c.user, c.password)
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func main() {
	fmt.Println("-------------Keychain-------------")
	k := CreateKeychain()
	fmt.Print("Insert 🔑: ")
	b := readLine()
	secret, err := secret.CreateSecret(b)
	if err != nil {
		log.Fatal(err)
	}

	if k.userAlreadyCreated() {
		k.setPassword(secret)
	} else {
		err := k.login(secret)

		if err != nil {
			color.Red.Println("Incorrect password!")
		} else {
			color.Green.Println("Success!")

			for i := range k.credentials {
				if i == len(k.credentials)-1 {
					color.BgYellow.Println(k.credentials[i].print(i))
				} else {
					color.BgBlack.Println(k.credentials[i].print(i))
				}
			}

			var currentIndex int = len(k.credentials) - 1

			ansi.CursorHide()

			fmt.Print("d: delete")

			ansi.CursorPreviousLine(0)

			for {
				value, arrowKey, err := keyboard.GetSingleKey()
				if err != nil {
					log.Fatal(err)
				}

				if value == rune('q') {
					ansi.CursorShow()
					os.Exit(0)
				}

				didMove := false

				if isNumber(value) {
					if indexToMove := int(value - '0'); indexToMove <= len(k.credentials)-1 {
						if linesToMoveDown := indexToMove - currentIndex; linesToMoveDown != 0 {
							ansi.CursorHorizontalAbsolute(0)
							color.BgBlack.Print(k.credentials[currentIndex].print(currentIndex))
							if linesToMoveDown > 0 {
								if newIndex := currentIndex + linesToMoveDown; newIndex <= len(k.credentials)-1 {
									ansi.CursorNextLine(linesToMoveDown)
									didMove = true
									currentIndex = newIndex
								}
							} else if linesToMoveDown != 0 {
								linesToMoveUp := -linesToMoveDown
								if newIndex := currentIndex - linesToMoveUp; newIndex >= 0 {
									ansi.CursorPreviousLine(linesToMoveUp)
									didMove = true
									currentIndex = newIndex
								}
							}
						}
					}
				} else if arrowKey == keyboard.KeyArrowUp {
					if currentIndex > 0 {
						ansi.CursorHorizontalAbsolute(0)
						color.BgBlack.Print(k.credentials[currentIndex].print(currentIndex))
						ansi.CursorPreviousLine(0)
						currentIndex--
						didMove = true
					}
				} else if arrowKey == keyboard.KeyArrowDown {
					if currentIndex < len(k.credentials)-1 {
						ansi.CursorHorizontalAbsolute(0)
						color.BgBlack.Print(k.credentials[currentIndex].print(currentIndex))
						ansi.CursorNextLine(0)
						currentIndex++
						didMove = true
					}
				}

				if didMove {
					color.BgYellow.Print(k.credentials[currentIndex].print(currentIndex))
				}
			}
		}
	}
}
