package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

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
	fmt.Println("Accepted connection from: ", conn.RemoteAddr())

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if (err != nil) {
		fmt.Println("error receiving data:", err)
	}

	req := string(buf[:n])

	// splits the request 
	parts := strings.Split(req, "\r\n")

	// 
	sec := strings.Split(parts[0], " ")

	if (sec[1] == "/") {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}