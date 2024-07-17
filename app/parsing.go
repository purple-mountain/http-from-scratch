package main

import (
	"fmt"
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

func ParseResponse(req string) string {
	requestLine := ParseRequestLine(req)
	httpVersion := 1.1

	var httpMethod, url, body, header string
	var statusCode, contentLength int

	unpack(strings.Split(requestLine, " "), &httpMethod, &url)

	urlSegments := strings.Split(url, "/")

	statusCode = 200
	if url == "/" {
		return fmt.Sprintf("HTTP/%1.1f %d %s\r\n\r\n", 1.1, statusCode, statusMsg[statusCode])
	}

	// TODO: strings.HasPrefix refactor
	if urlSegments[1] == "echo" && len(urlSegments) == 3 {
		body = urlSegments[2]
	} else if url == "/user-agent" {
		body = ParseUserAgent(req)
	} else {
		statusCode = 404
	}

	contentLength = len(body)
	header = fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n", contentLength)
	statusLine := fmt.Sprintf("HTTP/%1.1f %d %s\r\n", httpVersion, statusCode, statusMsg[statusCode])

	res := statusLine + header + "\r\n" + body

	return res
}
