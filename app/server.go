package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const HTTP_VERSION = 1.1

var statusMsg = map[int]string{
	200: "OK",
	404: "Not Found",
}

func unpack(s []string, vars ...*string) {
	for i := range vars {
		*vars[i] = s[i]
	}
}

func getRequestLine(req string) string {
	firstLineEndIDx := strings.Index(req, "\r\n")
	return req[:firstLineEndIDx]
}

func getUserAgent(req string) string {
	userAgentIdx := strings.Index(req, "User-Agent: ") + len("User-Agent: ")
	userAgentEndIdx := strings.Index(req[userAgentIdx:], "\r\n") + userAgentIdx

	return req[userAgentIdx:userAgentEndIdx]
}

func getResponse(req string) string {
	requestLine := getRequestLine(req)

	var httpMethod, url, body, header string
	var statusCode, contentLength int

	unpack(strings.Split(requestLine, " "), &httpMethod, &url)

	urlSegments := strings.Split(url, "/")

	statusCode = 200
	if url == "/" {
		return fmt.Sprintf("HTTP/%1.1f %d %s\r\n\r\n", HTTP_VERSION, statusCode, statusMsg[statusCode])
	}

	if urlSegments[1] == "echo" && len(urlSegments) == 3 {
		body = urlSegments[2]
	} else if url == "/user-agent" {
		body = getUserAgent(req)
	} else {
		statusCode = 404
	}

	contentLength = len(body)
	header = fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n", contentLength)
	statusLine := fmt.Sprintf("HTTP/%1.1f %d %s\r\n", HTTP_VERSION, statusCode, statusMsg[statusCode])

	res := statusLine + header + "\r\n" + body

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
