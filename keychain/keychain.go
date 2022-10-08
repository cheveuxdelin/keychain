package keychain

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cheveuxdelin/keychain/crypt"
	"github.com/cheveuxdelin/keychain/secret"
	"github.com/cheveuxdelin/keychain/utils"
	"github.com/eiannone/keyboard"
	"github.com/gookit/color"
	"github.com/k0kubun/go-ansi"
)

type settings struct {
	filename      string
	wordDelimiter byte
	lineDelimiter byte
}

type Keychain struct {
	credentials []credential
	secret      secret.Secret
	settings    settings
}

func (k Keychain) credentialsToBytes() []byte {
	var b bytes.Buffer
	for i := range k.credentials {
		b.WriteString(k.credentials[i].user)
		b.WriteByte(k.settings.wordDelimiter)
		b.WriteString(k.credentials[i].password)
		b.WriteByte(k.settings.lineDelimiter)
	}
	return b.Bytes()
}

var longestCredential = 0

const TERMINAL_MINIMUM_WIDTH int = 30

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
		var newCredentials credential = credential{
			user:     user.String(),
			password: password.String(),
		}
		longestCredential = utils.Max(longestCredential, newCredentials.Length())
		k.credentials = append(k.credentials, newCredentials)
		user.Reset()
		password.Reset()
	}
	return nil
}

func (k *Keychain) createCredential(user string, password string) {
	k.credentials = append(k.credentials, credential{user: user, password: password})
	k.save()
}

func (k *Keychain) setPassword(secret secret.Secret) {
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

func (k *Keychain) userAlreadyCreated() bool {
	_, err := os.Stat(k.settings.filename)
	return errors.Is(err, os.ErrNotExist)
}

func startKeyboard() {
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
}
func closeKeyboard() {
	if err := keyboard.Close(); err != nil {
		log.Fatal(err)
	}
}
func clearConsole() {
	fmt.Print("\033c")
}

// FLOW
func (k *Keychain) Start() {
	startKeyboard()
	clearConsole()
	k.auth()
	k.run()
}
func (k *Keychain) auth() {
	fmt.Print("Insert 🔑: ")
	secret := utils.GetEnteredPassword()
	if k.userAlreadyCreated() {
		k.setPassword(secret)
	} else if err := k.login(secret); err != nil {
		color.Red.Println("incorrect password 🚫")
		os.Exit(1)
	}
}
func (k *Keychain) run() {
	var safe bool = true
	for {
		clearConsole()
		fmt.Println("Keychain v0.3.0")
		PrintHeaders()
		var currentIndex int = len(k.credentials) - 1
		ansi.CursorHide()
		for i := range k.credentials {
			if i == len(k.credentials)-1 {
				if !safe {
					color.BgYellow.Println(k.credentials[i].Print(i))
				} else {
					color.BgYellow.Println(k.credentials[i].PrintSafe(i))
				}
			} else {
				if !safe {
					color.BgBlack.Println(k.credentials[i].Print(i))
				} else {
					color.BgBlack.Println(k.credentials[i].PrintSafe(i))
				}
			}
		}
		fmt.Print("a: add   d: delete   q: quit")
		ansi.CursorPreviousLine(0)

		for {
			value, arrowKey, err := keyboard.GetKey()
			if err != nil {
				log.Fatal(err)
			}
			if value == rune('q') {
				ansi.CursorShow()
				keyboard.Close()
				os.Exit(0)
			}

			didMove := false

			if utils.IsNumber(value) {
				if indexToMove := int(value - '0'); indexToMove <= len(k.credentials)-1 {
					if linesToMoveDown := indexToMove - currentIndex; linesToMoveDown != 0 {
						ansi.CursorHorizontalAbsolute(0)
						if !safe {
							color.BgBlack.Print(k.credentials[currentIndex].Print(currentIndex))
						} else {
							color.BgBlack.Print(k.credentials[currentIndex].PrintSafe(currentIndex))
						}

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
					if !safe {
						color.BgBlack.Print(k.credentials[currentIndex].Print(currentIndex))
					} else {
						color.BgBlack.Print(k.credentials[currentIndex].PrintSafe(currentIndex))
					}
					ansi.CursorPreviousLine(0)
					currentIndex--
					didMove = true
				}
			} else if arrowKey == keyboard.KeyArrowDown {
				if currentIndex < len(k.credentials)-1 {
					ansi.CursorHorizontalAbsolute(0)
					if !safe {
						color.BgBlack.Print(k.credentials[currentIndex].Print(currentIndex))
					} else {
						color.BgBlack.Print(k.credentials[currentIndex].PrintSafe(currentIndex))
					}
					ansi.CursorNextLine(0)
					currentIndex++
					didMove = true
				}
			} else if value == 'a' {
				k.CreateCredential()
				break
			} else if value == 'd' && len(k.credentials) > 0 {
				k.DeleteCredential(currentIndex)
				break
			} else if value == 's' {
				safe = false
				break
			} else if value == 'h' {
				safe = true
				break
			}

			if didMove {
				if !safe {
					color.BgYellow.Print(k.credentials[currentIndex].Print(currentIndex))
				} else {
					color.BgYellow.Print(k.credentials[currentIndex].PrintSafe(currentIndex))
				}
			}
		}
	}
}

func CreateKeychain() (k *Keychain) {
	return &Keychain{
		settings: settings{
			filename:      ".credentials",
			wordDelimiter: ',',
			lineDelimiter: '\n',
		},
	}
}
