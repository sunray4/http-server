// current branch is master, and the remote location to personal repo is myorigin

package main

import (
	"bytes"
	"compress/gzip"
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
	// defer conn.Close()
	var res string 

	for {
		buf := make([]byte, 1024) // size of http request currently limited to 1024 bytes
		n, err := conn.Read(buf)

		if (err != nil) {
			fmt.Println("error receiving data:", err)
		}

		req := string(buf[:n])

		// splits the request in the request line, content in headers and body
		parts := strings.Split(req, "\r\n")

		// split the request line into the http method, request target and html version
		sec := strings.Split(parts[0], " ") 

		switch sec[0] {
			case "GET":
				res = handleGet(parts, sec)
			case "POST":
				res = handlePost(parts, sec)
		}

		close := checkConnClose(parts)
		if close {
			res = addCloseHeader(res)
			conn.Write([]byte(res))
			conn.Close()
			break
		} else {
			conn.Write([]byte(res))
		}

		
	}
	
}

func handleGet(parts []string, sec []string) (string) {
	var res string
	if sec[1] == "/" {
		res = "HTTP/1.1 200 OK\r\n\r\n" 
	} else if strings.Contains(sec[1], "/files/") {
		filename := strings.Split(sec[1], "/")[2]
		// path, err := fileSearch(".", filename)
		path := fmt.Sprintf("/tmp/data/codecrafters.io/http-server-tester/%s", filename)
		_, err := os.Stat(path)
		if err != nil {
			fmt.Print("error:", err)
			res = "HTTP/1.1 404 Not Found\r\n\r\n"
		} else {
			data, err := os.ReadFile(path)
			if (err != nil) {
				res = "HTTP/1.1 404 Not Found\r\n\r\n"
			} else {
				res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
			}
		}
	} else if sec[1] == "/user-agent" {
		var resBody string
		for i := 1; i < len(parts); i++ {
			if strings.Contains(strings.ToLower(parts[i]), "user-agent") {
				resBody = strings.TrimSpace(strings.Split(parts[i], ":")[1])
				break;
			}
		}
		res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(resBody), resBody)
	} else if strings.Contains(sec[1], "/echo/") {
		res = echoCompression(parts, sec)
	} else {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	return res
}

func handlePost(parts []string, sec []string) (string) {
	var res string
	if strings.Contains(sec[1], "/files") {
		filename := strings.Split(sec[1], "/")[2]
		filepath := fmt.Sprintf("/tmp/data/codecrafters.io/http-server-tester/%s", filename)
		content := []byte(parts[len(parts) - 1])

		err := os.WriteFile(filepath, content, 0664)

		if err!= nil {
			fmt.Println("error writing file:", err)
			res = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
		} else {
			res = "HTTP/1.1 201 Created\r\n\r\n"
		}

	} else {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	return res
}

func echoCompression(parts []string, sec []string) (string) {
	var encodings []string
	var res string
	supportedCompressions := make(map[string]struct{})
	supportedCompressions["gzip"] = struct{}{} // this server currently only supports gzip
	for i := range parts {
		if strings.Contains(strings.ToLower(parts[i]), "accept-encoding") {
			encodings = strings.Split(strings.Split(parts[i], ":")[1], ",")
			for i := range encodings {
				encodings[i] = strings.TrimSpace(encodings[i])
			}
			break
		}
	}

	if len(encodings) == 0 || encodings[0] == "invalid-encoding" {
		res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(sec[1][6:]), sec[1][6:])
	} else {
		for i := range encodings {
			_, exists := supportedCompressions[encodings[i]]
			if exists && encodings[i] == "gzip" {
				gzCont, cont := gzipWrite(sec[1][6:])
				if cont != "" {
					res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", encodings[i], len(sec[1][6:]), cont)
				} else {
					res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", encodings[i], len(gzCont), gzCont)
				}
				
			}
		}

		if res == "" {
			res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(sec[1][6:]), sec[1][6:])
		}
		
	}
	return res
}

func gzipWrite(content string) ([]byte, string) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(content))
	if err != nil {
		fmt.Println("error compressing with gzip:", err)
		return nil, content
	}

	if err := zw.Close(); err != nil {
		fmt.Println("error closing gzip:", err)
	}

	if buf.Len() != 0 {
		return buf.Bytes(), ""
	} else {
		return nil, content
	}
}

func checkConnClose(parts []string) (bool) {
	for i := 1; i < len(parts); i++ {
		if strings.Contains(strings.ToLower(parts[i]), "connection") {
			split := strings.SplitN(parts[i], ":", 2)
			var result string
			if len(split) > 1 {
				result = strings.ToLower(strings.TrimSpace(split[1]))
				if result == "close" {
					return true
				}
			} 
		}
	}
	return false
}

func addCloseHeader(res string) (string) {
	parts := strings.SplitN(res, "\r\n", 2)
	if len(parts) > 1 {
		newRes := parts[0] + "\r\n" + "Connection: close\r\n" + parts[1]
		return newRes
	} else {
		return res
	}
	

	
}

// func fileSearch(root string, target string) (string, error) {
// 	var found string
// 	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if err != nil  {
// 			return err
// 		}
// 		if !info.IsDir() && info.Name() == target {
// 			found = path
// 			// return filepath.SkipDir
// 		}
// 		if !info.IsDir() {
// 			fmt.Println(info.Name())
// 		}
// 		return err
// 	})
// 	if err != nil {
// 		return "", err
// 	} else if found == "" {
// 		return "", fmt.Errorf("file not found")
// 	} else {
// 		return found, nil
// 	}
// }