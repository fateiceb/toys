package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	str := sha256.Sum256([]byte("world"))
	fmt.Printf("%x", str)
}
