package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// TODO: Maybe I can use just raw bytes instead of strings

const HTTP_VERSION = 1.1

var statusMsg = map[int]string{
	200: "OK",
	404: "Not Found",
}

func getResponse(req string) string {
	firstLineEndIDx := strings.Index(req, "\r\n")
	requestLine := strings.Split(req[:firstLineEndIDx], " ")

	var res string
	var statusCode int

	if requestLine[1] == "/" {
		statusCode = 200
	} else {
		statusCode = 404
	}

	res = fmt.Sprintf("HTTP/%1.1f %d %s\r\n\r\n", HTTP_VERSION, statusCode, statusMsg[statusCode])

	return res
}

func handleClient(conn net.Conn) error {
	reqBuffer := make([]byte, 1024)

	_, err := conn.Read(reqBuffer)
	if err != nil {
		return err
	}

	res := getResponse(string(reqBuffer))

	conn.Write([]byte(res))

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

	err = handleClient(conn)
	if err != nil {
		fmt.Println("Error reading the response: ", err.Error())
		os.Exit(1)
	}
}
