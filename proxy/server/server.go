package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	log.Println("Starting HTTP proxy server on :8080")
	// HTTP CONNECT does not accept paths
	// so we'll just use the mux
	err := http.ListenAndServe(":8080", handler())
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s from %s",
			r.Method,
			r.Host,
			r.RemoteAddr,
		)
		if r.Method != http.MethodConnect {
			http.Error(w, "only HTTP CONNECT allowed", http.StatusMethodNotAllowed)
			return
		}
		handle(w, r)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	dest := r.Host
	// Establish TCP connection
	serverConn, err := net.Dial("tcp", dest)
	if err != nil {
		http.Error(w, "Unable to connect to destination", http.StatusServiceUnavailable)
		return
	}
	defer serverConn.Close()
	log.Printf("Proxy connected to %s", dest)
	log.Printf("Connection details: %s to %s",
		serverConn.LocalAddr(),
		serverConn.RemoteAddr(),
	)

	// Upgrade connection to a tunnel
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Unable to hijack connection", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Unable to hijack connection", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// Return 200 to client
	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	log.Println("200 returned to client")

	// Proxy data
	go func() {
		_, err = io.Copy(serverConn, clientConn) // Client to Server
		if err != nil {
			log.Printf("Error copying data from client to server: %v", err)
		}
	}()
	_, err = io.Copy(clientConn, serverConn) // Server to Client
	if err != nil {
		log.Printf("Error copying data from server to client: %v", err)
	}
}
