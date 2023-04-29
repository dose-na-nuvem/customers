package main

import (
	"fmt"
)

func main() {
	msg := message()
	fmt.Println(msg)
}

func message() string {
	return "hello world"
}
