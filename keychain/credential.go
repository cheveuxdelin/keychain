package keychain

import (
	"fmt"
	"strings"
)

type credential struct {
	user     string
	password string
}

func (c *credential) Print(indexNumber int) string {
	return fmt.Sprintf("[%d] %-20s | %-20s", indexNumber, c.user, c.password)
}

func (c *credential) PrintSafe(indexNumber int) string {
	return fmt.Sprintf("[%d] %-20s | %-20s", indexNumber, c.user, strings.Repeat("•", len(c.password)))
}
