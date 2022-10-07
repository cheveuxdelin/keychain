package main

import "fmt"

type credential struct {
	user     string
	password string
}

func (c *credential) print(indexNumber int) string {
	return fmt.Sprintf("[%d] %-20s | %-20s", indexNumber, c.user, c.password)
}
