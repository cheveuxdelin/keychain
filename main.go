package main

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/cheveuxdelin/keychain/crypt"
	"github.com/cheveuxdelin/keychain/secret"
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

func main() {

}
