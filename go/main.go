package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
)

const (
	dnsServerAddr = "1.1.1.1"
	dnsServerPort = "853"
)

var config = tls.Config{
	ServerName:         "1.1.1.1",
	InsecureSkipVerify: true,
}

func copy(src net.Conn, dest net.Conn) error {
	_, err := io.Copy(dest, src)
	if err != nil {
		fmt.Println("Error copying data to server:", err)
	}
	return err
}

func startProxyTCP(listener net.Listener) error {

	fmt.Println("DNS proxy listening on port 853")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Connecting to DNS server...")
		serverTLSConn, err := tls.Dial("tcp", dnsServerAddr+":"+dnsServerPort, &config)
		if err != nil {
			fmt.Println("Error connecting to DNS server:", err)
			conn.Close()
			continue
		}

		fmt.Println("Connected to DNS server.")
		go copy(conn, serverTLSConn)
		go copy(serverTLSConn, conn)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":53")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	startProxyTCP(listener)
}
