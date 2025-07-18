package main

import (
	"fmt"
	"net/http"
)

func main() {

	// l, err := net.Listen("tcp", "0.0.0.0:4221")
	// if err != nil {
	// 	fmt.Println("Failed to bind to port 4221")
	// 	os.Exit(1)
	// }
	
	// conn, err := l.Accept()
	// if err != nil {
	// 	fmt.Println("Error accepting connection: ", err.Error())
	// 	os.Exit(1)
	// }
	// fmt.Println("Accepted connection from: ", conn.RemoteAddr())

	// conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	http.HandleFunc("/", handleRoot)
	
	http.ListenAndServe(":4221", nil)

}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL)
	if req.URL.Path == "/" {
		fmt.Fprintf(w, "HTTP/1.1 200 OK\r\n\r\n")
	} else {
		fmt.Fprintf(w,"HTTP/1.1 400 Not Found\r\n\r\n")
	}
	
}