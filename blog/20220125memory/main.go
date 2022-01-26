package main

import (
	"bytes"
	"io"
)

func algOne(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
	input := bytes.NewBuffer(data)
	size := len(find)
	buf := make([]byte, size)
	end := size - 1
	if n, err := io.ReadFull(input, buf[:end]); err != nil {
		output.Write(buf[:n])
		return
	}
	for {
		if _, err := io.ReadFull(input, buf[:end]); err != nil {
			output.Write(buf[:end])
			return
		}
		if bytes.Compare(buf, find) == 0 {
			output.Write(repl)
			if n, err := io.ReadFull(input, buf[:end]); err != nil {
				output.Write(buf[:n])
				return
			}
			continue
		}
		output.WriteByte(buf[0])
		copy(buf, buf[1:])
	}
}
func main() {

}
