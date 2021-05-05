package main

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func Test_generateId(t *testing.T) {
	hash := sha256.New()
	hash.Write([]byte("1"))
	want := hash.Sum(nil)
	got := generateId("1")
	t.Run("Id want == Id got", func(t *testing.T) {

		if !bytes.Equal(want, got) {
			t.Fatal("got != want", want, got)
		}
	})
}
