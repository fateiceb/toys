package main

import (
	"crypto/elliptic"
	crand "crypto/rand"
	"fmt"
)

func main()  {
	key, b, b2, err := elliptic.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(key,b,b2,len(key))
}
