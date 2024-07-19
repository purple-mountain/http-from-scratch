package main

import (
	"fmt"
	"os"
	"strings"
)

func ParseRequestLine(req string) string {
	firstLineEndIDx := strings.Index(req, "\r\n")
	return req[:firstLineEndIDx]
}

func ParseUserAgent(req string) string {
	userAgentIdx := strings.Index(req, "User-Agent: ") + len("User-Agent: ")
	userAgentEndIdx := strings.Index(req[userAgentIdx:], "\r\n") + userAgentIdx

	return req[userAgentIdx:userAgentEndIdx]
}

func ParseResponse(req, filesPath string) string {
	requestLine := ParseRequestLine(req)
	httpVersion := 1.1
	contentType := "text/plain"

	var httpMethod, url, body, header string
	var statusCode, contentLength int

	unpack(strings.Split(requestLine, " "), &httpMethod, &url)

	urlSegments := strings.Split(url, "/")

	statusCode = 200
	if url == "/" {
		return fmt.Sprintf("HTTP/%1.1f %d %s\r\n\r\n", httpVersion, statusCode, statusMsg[statusCode])
	}

	// TODO: strings.HasPrefix refactor
	if urlSegments[1] == "echo" && len(urlSegments) == 3 {
		body = urlSegments[2]
	} else if strings.HasPrefix(url, "/files/") {
		fileName := urlSegments[2]
		data, err := os.ReadFile(filesPath + fileName)
		if err != nil {
			statusCode = 404
			return fmt.Sprintf("HTTP/%1.1f %d %s\r\n\r\n", httpVersion, statusCode, statusMsg[statusCode])
		}

		body = string(data)
		contentType = "application/octet-stream"
	} else if url == "/user-agent" {
		body = ParseUserAgent(req)
	} else {
		statusCode = 404
	}

	contentLength = len(body)
	header = fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n", contentType, contentLength)
	statusLine := fmt.Sprintf("HTTP/%1.1f %d %s\r\n", httpVersion, statusCode, statusMsg[statusCode])

	res := statusLine + header + "\r\n" + body

	return res
}
