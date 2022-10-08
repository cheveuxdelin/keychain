package main

import (
	"github.com/cheveuxdelin/keychain/keychain"
)

func main() {
	keychain.CreateKeychain().Start()
}
