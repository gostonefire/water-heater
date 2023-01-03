package main

import (
	"fmt"
	"net"
)

func main() {
	cmd := "vout:00000000.0"
	address := "192.168.1.138:5000"

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("error while connecting to server: %s\n", err)
		return
	}
	defer func(conn net.Conn) { _ = conn.Close() }(conn)

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		fmt.Printf("error while writing to server: %s\n", err)
		return
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("error while reading from server: %s\n", err)
		return
	}

	fmt.Println(n)
	fmt.Println(string(buf))

}
