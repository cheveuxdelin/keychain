package keychain

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gookit/color"
)

type credential struct {
	user     string
	password string
}

var numberColumnWidth int = 5
var userColumnWidth int = 15
var passwordColumnWith int = 20

func PrintHeaders() {
	color.Set(color.BgGreen, color.FgCyan)
	spaceForNumber := strings.Repeat(" ", numberColumnWidth/2)
	number := spaceForNumber + "#" + spaceForNumber
	fmt.Print(number)
	fmt.Print("|")
	fmt.Printf("%-*s", userColumnWidth, "user")
	fmt.Print("|")
	fmt.Printf("%-*s", passwordColumnWith, "password")
	fmt.Println()
	color.Reset()
}

func (c *credential) Print(indexNumber int) string {
	spaceForNumber := strings.Repeat(" ", numberColumnWidth/2-1)
	paddedNumber := fmt.Sprintf("%s[%d]%s", spaceForNumber, indexNumber, spaceForNumber)
	paddedUser := fmt.Sprintf("%-*s", userColumnWidth, c.user)
	paddedPassword := fmt.Sprintf("%-*s", passwordColumnWith, c.password)
	return fmt.Sprint(strings.Join([]string{paddedNumber, paddedUser, paddedPassword}, "|"))
}

func (c *credential) PrintSafe(indexNumber int) string {
	spaceForNumber := strings.Repeat(" ", numberColumnWidth/2-1)
	paddedNumber := fmt.Sprintf("%s[%d]%s", spaceForNumber, indexNumber, spaceForNumber)

	paddedUser := fmt.Sprintf("%-*s", userColumnWidth, c.user)
	paddedPassword := fmt.Sprintf("%-*s", passwordColumnWith, strings.Repeat("*", len(c.password)))
	return fmt.Sprint(strings.Join([]string{paddedNumber, paddedUser, paddedPassword}, "|"))
}

func (c *credential) Length() int {
	return len(c.user) + len(c.password)
}

func (c *credential) Copy() {
	clipboard.WriteAll(c.user + "|" + c.password)
}
