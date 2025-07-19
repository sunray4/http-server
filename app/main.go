// current branch is master, and the remote location to personal repo is myorigin

package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Accepted connection from: ", conn.RemoteAddr())

		go handleConnection(conn)
	}

	
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if (err != nil) {
		fmt.Println("error receiving data:", err)
	}

	req := string(buf[:n])

	// splits the request in the request line, content in headers and body
	parts := strings.Split(req, "\r\n")

	// split the request line into the http method, request target and html version
	sec := strings.Split(parts[0], " ") 

	if (sec[1] == "/") {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n")) // respond to the request
	} else if (strings.Contains(sec[1], "/files/")) {
		filename := strings.Split(sec[1], "/")[2]
		// path, err := fileSearch("/private/tmp", filename)
		path := fmt.Sprintf("/tmp/data/codecrafters.io/http-server-tester/%s", filename)
		_, err := os.Stat(path)
		if (err != nil) {
			fmt.Print("error:", err)
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			data, err := os.ReadFile(path)
			if (err != nil) {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
			res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
			conn.Write([]byte(res) ) // respond to the request
		}
		
	} else if (sec[1] == "/user-agent") {
		id := 0
		for i := 0; i < len(parts); i++ {
			if (strings.Contains(parts[i], "User-Agent")) {
				id = i;
				break;
			}
		}
		resBody := strings.Split(parts[id], " ")[1]
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(resBody), resBody)
		conn.Write([]byte(res) )
	} else if (strings.Contains(sec[1], "/echo/")) {
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(sec[1][6:]), sec[1][6:])
		conn.Write([]byte(res) ) // respond to the request
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func fileSearch(root string, target string) (string, error) {
	var found string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil  {
			return err
		}
		if !info.IsDir() && info.Name() == target {
			found = path
			// return filepath.SkipDir
		}
		if !info.IsDir() {
			fmt.Println(info.Name())
		}
		return err
	})
	if err != nil {
		return "", err
	} else if found == "" {
		return "", fmt.Errorf("file not found")
	} else {
		return found, nil
	}
}