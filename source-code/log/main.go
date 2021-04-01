package main

import (
	"bytes"
	"log"
	"os"
)

var (
	buffer bytes.Buffer
)

func main() {
	file, err := os.OpenFile("mainlog", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "main：", log.Ldate|log.Ltime|log.Lshortfile)
	// logger.SetOutput(&buffer)
	logger.Print("hahahahaha")
}
