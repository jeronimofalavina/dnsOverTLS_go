package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	dnsServerAddr    = "1.1.1.1"
	dnsServerPortTLS = "853"
	dnsServerPortUDP = "53"
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

func startProxyUDP(udpConn net.PacketConn) error {
	fmt.Println("DNS proxy listening on UDP port 53")

	for {
		buffer := make([]byte, 512)
		n, addr, err := udpConn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading request:", err)
			continue
		}

		fmt.Printf("Received DNS request from %s\n", addr.String())

		conn, err := net.Dial("udp", dnsServerAddr+":"+dnsServerPortUDP)
		if err != nil {
			fmt.Println("Error connecting to DNS server:", err)
			continue
		}

		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error sending request to DNS server:", err)
			continue
		}

		buffer = make([]byte, 512)
		n, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading response from DNS server:", err)
			continue
		}

		fmt.Printf("Received DNS response from %s udp\n", dnsServerAddr)

		_, err = udpConn.WriteTo(buffer[:n], addr)
		if err != nil {
			fmt.Println("Error sending response to client:", err)
			continue
		}
		conn.Close()
	}

}

func startProxyTCP(listener net.Listener) error {

	fmt.Println("DNS proxy listening on TCP port 53")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Connecting to DNS server...")
		serverTLSConn, err := tls.Dial("tcp", dnsServerAddr+":"+dnsServerPortTLS, &config)
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
	log.Println("Hello world!")
	udpListener, err := net.ListenPacket("udp", ":5333")
	if err != nil {
		panic(err)
	}

	tcpListener, err := net.Listen("tcp", ":5333")
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()

	go startProxyUDP(udpListener)
	go startProxyTCP(tcpListener)

	fmt.Scanln()
}
