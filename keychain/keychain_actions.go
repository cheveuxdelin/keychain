package keychain

import (
	"fmt"

	"github.com/cheveuxdelin/keychain/utils"
	"github.com/k0kubun/go-ansi"
)

func (k *Keychain) CreateCredential() {
	closeKeyboard()
	clearConsole()
	ansi.CursorShow()
	fmt.Print("user: ")
	user := utils.ReadString()
	fmt.Print("password: ")
	password := utils.ReadSafeBytes()
	k.createCredential(user, string(password))
	k.save()
	fmt.Println(user, password)
	startKeyboard()
}

func (k *Keychain) DeleteCredential(indexToDelete int) {
	k.credentials = append(k.credentials[:indexToDelete], k.credentials[indexToDelete+1:]...)
	k.save()
}
