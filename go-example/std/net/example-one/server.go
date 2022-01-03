package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {

	ln, err := net.Listen("tcp", "127.0.0.1:10001")
	defer ln.Close()
	if err != nil {
		fmt.Errorf("err%s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Errorf("err%s", err)
		}
		go func(c net.Conn) {
			io.Copy(c, c)
			rd := bufio.NewReader(c)
			result, err := rd.ReadString('\n')
			if err != nil {
				fmt.Errorf("err%s", err)
			}
			fmt.Println(result)
			c.Close()
		}(conn)
	}

}
