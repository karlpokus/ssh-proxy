package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

// io pipes:
//
// ssh client stdout is piped to our stdin
// our stdout is piped to ssh client stdin

func main() {
	fpath := "./cmd.log"
	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// nowhere to log yet
		return
	}
	defer f.Close()
	flog := log.New(f, "", log.LstdFlags)

	sshHost := os.Args[1]
	sshPort := os.Args[2]
	proxyHost := "localhost"
	proxyPort := "8080"
	// proxyUsername := "your_username"
	// proxyPassword := "your_password"
	flog.Printf("config: sshHost: %s sshPort: %s proxyHost: %s proxyPort: %s",
		sshHost,
		sshPort,
		proxyHost,
		proxyPort,
	)

	// Create the connection to the proxy server
	proxyConn, err := net.Dial("tcp", net.JoinHostPort(proxyHost, proxyPort))
	if err != nil {
		flog.Fatalf("dial proxy server fail: %v", err)
	}
	defer proxyConn.Close()

	// Prepare the HTTP CONNECT request
	//auth := base64.StdEncoding.EncodeToString([]byte(proxyUsername + ":" + proxyPassword))
	//connectRequest := fmt.Sprintf("CONNECT %s:%s HTTP/1.1\r\nHost: %s:%s\r\nProxy-Authorization: Basic %s\r\n\r\n",
	connectRequest := fmt.Sprintf("CONNECT %s:%s HTTP/1.1\r\nHost: %s:%s\r\n\r\n",
		sshHost, sshPort, sshHost, sshPort)

	// Send CONNECT
	_, err = proxyConn.Write([]byte(connectRequest))
	if err != nil {
		flog.Fatalf("send connect fail: %v", err)
	}

	// Read response
	resp, err := http.ReadResponse(bufio.NewReader(proxyConn), nil)
	if err != nil || resp.StatusCode != 200 {
		flog.Fatalf("read proxy respons fail: %v status: %s", err, resp.Status)
	}
	flog.Println("Tunnel established")
	defer resp.Body.Close()

	// forwarding traffic
	go func() {
		_, err = io.Copy(proxyConn, os.Stdin)
		if err != nil {
			flog.Printf("copy stdin to proxy fail: %v", err)
		}
	}()

	_, err = io.Copy(os.Stdout, proxyConn)
	if err != nil {
		flog.Fatalf("copy proxy to stdout fail: %v", err)
	}
}
