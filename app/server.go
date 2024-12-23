package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var statusMsg = map[int]string{
	200: "OK",
	404: "Not Found",
	201: "Created",
}

type Server struct {
	listener  net.Listener
	conn      chan net.Conn
	shutdown  chan struct{}
	error     chan error
	filesPath string
}

func (s *Server) HandleClient(conn net.Conn) {
	defer conn.Close()
	reqBuffer := make([]byte, 1024)

	_, err := conn.Read(reqBuffer)
	if err != nil {
		s.error <- err
		return
	}

	res, err := ParseResponse(string(reqBuffer), s.filesPath)
	if err != nil {
		s.error <- err
	}

	conn.Write([]byte(res))
}

func (s *Server) AcceptConnections() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		s.conn <- conn
	}
}

func (s *Server) HandleConnections() {
	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.conn:
			go s.HandleClient(conn)
			// TODO: add error handling
		}
	}
}

func (s *Server) Stop() {
	close(s.conn)
	close(s.shutdown)
}

func NewServer(port, filesPath string) (*Server, error) {
	l, err := net.Listen("tcp", "0.0.0.0"+port)
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}

	return &Server{
		listener:  l,
		conn:      make(chan net.Conn),
		shutdown:  make(chan struct{}),
		error:     make(chan error),
		filesPath: filesPath,
	}, nil
}

func main() {
	filesPath := flag.String("directory", "/tmp/", "File storage location")
	flag.Parse()

	server, err := NewServer(":4221", *filesPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go server.AcceptConnections()
	go server.HandleConnections()

	time.Sleep(60 * time.Second)
	server.Stop()

	fmt.Println("Server was shutdown")
}
