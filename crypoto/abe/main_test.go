package main

import (
	"fmt"
	"log"
	"testing"
)

func Test_byteAdd(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"addbyte"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := NewPrivKey()
			log.Println("privatekey:", key)
			fmt.Printf("%d\n", key)
			total := make([]byte, 32)
			log.Println(total)
			final := byteAdd(total, key)
			log.Fatal("byteadd:", final)
		})
	}
}
