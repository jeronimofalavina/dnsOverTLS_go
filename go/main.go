package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
)

func main() {
	// Specify the address of the DNS server to forward requests to
	dnsServerAddr := "1.1.1.1"
	dnsServerPort := "853"

	// Load TLS configuration
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "1.1.1.1",
		//	InsecureSkipVerify: true,
	}
	// Listen for incoming DNS requests on port 853
	listener, err := net.Listen("tcp", ":5333")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("DNS proxy listening on port 853")

	// Loop forever, handling incoming requests
	for {
		// Accept the incoming DNS request
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Connecting to DNS server...")
		// Attempt to establish a TLS connection to the DNS server
		serverTLSConn, err := tls.Dial("tcp", dnsServerAddr+":"+dnsServerPort, &config)
		if err != nil {
			fmt.Println("Error connecting to DNS server:", err)
			conn.Close()
			continue
		}

		fmt.Println("Connected to DNS server.")

		go func() {
			_, err := io.Copy(serverTLSConn, conn)
			if err != nil {
				fmt.Println("Error copying data to server:", err)
			}

		}()

		go func() {
			_, err := io.Copy(conn, serverTLSConn)
			if err != nil {
				fmt.Println("Error copying data to client:", err)
			}
			conn.Close()
			serverTLSConn.Close()
		}()

	}
}
