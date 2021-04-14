package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// // str := sha256.Sum256([]byte("world"))
	// message := "hellodsadsadsads2ad"
	// hash := sha256.New()
	// hash.Write([]byte(message))
	// sign := hash.Sum([]byte(message))
	// fmt.Printf("%s-----%x-----%d", string(sign[0:len(sign)-32]), sign[len(sign)-32:len(sign)], len(sign))
	// // fmt.Printf("%x%d", str, len(str))
	// os.IsPermission()
	// file, err := os.OpenFile("mykey.pem", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()
	// //生成rsa
	// var privkey *rsa.PrivateKey
	// if privkey, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
	// 	log.Fatal(err)
	// }
	// data := x509.MarshalPKCS1PrivateKey(privkey)
	// if err = pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: data}); err != nil {
	// 	log.Fatal(err)
	// }

	// hash := sha256.New()
	// hash.Write([]byte("111hello"))
	// // hash.Write([]byte("111"))
	// bytea := hash.Sum(nil)
	// log.Println(bytea, len(bytea))
	// // hash := sha256.New()
	// hahs := hmac.New(sha256.New, []byte("111"))
	// hahs.Write([]byte("hello"))
	// byteb := hahs.Sum(nil)
	// log.Println(byteb, len(byteb))
	// timenow := time.Now().Unix()

	// log.Println(timenow)
	// log.Println(timenowformat)
	rand.Seed(23)
	num := rand.Int31()
	fmt.Println()
}
