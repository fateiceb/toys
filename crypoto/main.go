package main

import (
	"crypto"
	"crypto/sha256"
)

func main() {
	hash := sha256.New()
	hash.Write([]byte("13131"))
	res := hash.Sum(nil)
	crypto.SHA3_384
}
