package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:10001")
	if err != nil {
		fmt.Errorf("%s", err)
	}
	_, err = io.WriteString(conn, "hello worrld\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil { //
		fmt.Errorf("%s", err)
	}
	fmt.Println(status)

}
