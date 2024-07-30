package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"slices"
	"strings"
)

func parseContentEncoding(req string) []string {
	idx := strings.Index(req, "Accept-Encoding: ") + len("Accept-Encoding: ")
	idxEnd := strings.Index(req[idx:], "\r\n") + idx

	return strings.Split(req[idx:idxEnd], ", ")
}

func ParseRequestLine(req string) string {
	firstLineEndIDx := strings.Index(req, "\r\n")

	return req[:firstLineEndIDx]
}

func ParseUserAgent(req string) string {
	userAgentIdx := strings.Index(req, "User-Agent: ") + len("User-Agent: ")
	userAgentEndIdx := strings.Index(req[userAgentIdx:], "\r\n") + userAgentIdx

	return req[userAgentIdx:userAgentEndIdx]
}

func ParseResponse(req, filesPath string) (string, error) {
	var httpMethod, url, body, header, contentEncoding string
	var contentLength int

	requestLine := ParseRequestLine(req)
	httpVersion := 1.1
	contentType := "text/plain"
	statusCode := 200
	acceptedEncodings := parseContentEncoding(req)

	unpack(strings.Split(requestLine, " "), &httpMethod, &url)

	urlSegments := strings.Split(url, "/")

	if url == "/" {
		return fmt.Sprintf("HTTP/%1.1f %d %s\r\n\r\n", httpVersion, statusCode, statusMsg[statusCode]), nil
	}

	if strings.HasPrefix(url, "/echo") {
		body = urlSegments[2]
	} else if strings.HasPrefix(url, "/files/") {
		fileName := urlSegments[2]

		if httpMethod == "POST" {
			statusCode = 201
			reqBody := req[strings.Index(req, "\r\n\r\n")+len("\r\n\r\n") : strings.Index(req, "\x00")]
			fmt.Println(reqBody)
			file, err := os.Create(filesPath + fileName)
			if err != nil {
				return "", err
			}
			_, err = file.WriteString(reqBody)
			if err != nil {
				return "", err
			}

			file.Sync()
			return fmt.Sprintf("HTTP/%1.1f %d %s\r\n\r\n", httpVersion, statusCode, statusMsg[statusCode]), nil
		}

		data, err := os.ReadFile(filesPath + fileName)
		if err != nil {
			statusCode = 404
			return fmt.Sprintf("HTTP/%1.1f %d %s\r\n\r\n", httpVersion, statusCode, statusMsg[statusCode]), nil
		}

		body = string(data)
		contentType = "application/octet-stream"
	} else if url == "/user-agent" {
		body = ParseUserAgent(req)
	} else {
		statusCode = 404
	}

	if slices.Contains(acceptedEncodings, "gzip") {
		var buffer bytes.Buffer
		gWriter := gzip.NewWriter(&buffer)
		gWriter.Write([]byte(body))
		gWriter.Close()
		body = buffer.String()
		contentEncoding = "Content-Encoding: gzip\r\n"
	}

	contentLength = len(body)
	header = fmt.Sprintf("Content-Type: %s\r\n%sContent-Length: %d\r\n", contentType, contentEncoding, contentLength)
	statusLine := fmt.Sprintf("HTTP/%1.1f %d %s\r\n", httpVersion, statusCode, statusMsg[statusCode])

	res := statusLine + header + "\r\n" + body

	return res, nil
}
