package keychain

func CreateTestKeychain() (k *Keychain) {
	return &Keychain{
		settings: settings{
			filename:      "test",
			wordDelimiter: ',',
			lineDelimiter: '\n',
		},
	}
}
