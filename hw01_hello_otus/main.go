package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	message := "Hello, OTUS!"
	reverseMessage := reverse.String(message)
	fmt.Println(reverseMessage)
}
