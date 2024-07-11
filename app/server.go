package main

import (
	"fmt"
	"net"
	"os"
)

const HTTP_VERSION = 1.1

func getResponse(req []byte) string {
	res := fmt.Sprintf("HTTP/%1.1f 200 OK\r\n\r\n", HTTP_VERSION)

	return res
}

func handleClient(conn net.Conn) error {
	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	res := getResponse(buffer)

	conn.Write([]byte(res))

	// GET / HTTP / 1.1

	// HTTP/1.1 200 OK\r\n\r\n

	return nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	handleClient(conn)

	// conn.Read()
}
