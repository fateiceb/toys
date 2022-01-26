package main

import (
	"bytes"
	"testing"
)

func BenchmarkAlgorighmOne(b *testing.B) {
	var output bytes.Buffer
	in := []byte("abcelvisaElvisabcelviseelvisaelvisaabeeeelvise l v i saa bb e l v i saa elviselvielviselvielvielviselvi1elvielviselvis")
	find := []byte("elvis")
	repl := []byte("Elvis")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output.Reset()
		algOne(in, find, repl, &output)
	}
}
