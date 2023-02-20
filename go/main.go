package main

import (
	"crypto/tls"
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
	ServerName: "1.1.1.1",
}

func startProxyUDP(udpConn net.PacketConn) error {
	log.Println("[UDP] - DNS proxy listening on UDP port 53")

	for {
		buffer := make([]byte, 512)
		n, addr, err := udpConn.ReadFrom(buffer)
		if err != nil {
			log.Println("[UDP] - Error reading request:", err)
			continue
		}

		log.Printf("[UDP] - Received DNS request from %s\n", addr.String())
		log.Println("[UDP] - Connecting to DNS server...")
		conn, err := net.Dial("udp", dnsServerAddr+":"+dnsServerPortUDP)
		if err != nil {
			log.Println("[UDP] - Error connecting to DNS server:", err)
			conn.Close()
			continue
		}

		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Println("[UDP] - Error sending request to DNS server:", err)
			conn.Close()
			continue
		}

		buffer = make([]byte, 512)
		n, err = conn.Read(buffer)
		if err != nil {
			log.Println("[UDP] - Error reading response from DNS server:", err)
			conn.Close()
			continue
		}

		// TODO: log the response from the dns server
		log.Printf("[UDP] - Received DNS response from %s \n", dnsServerAddr)

		_, err = udpConn.WriteTo(buffer[:n], addr)
		if err != nil {
			log.Println("[UDP] - Error sending response to client:", err)
			conn.Close()
			continue
		}
		conn.Close()
	}

}

func copy(src net.Conn, dest net.Conn) error {
	_, err := io.Copy(dest, src)
	if err != nil {
		log.Panicln("[TCP] - Error copying from message from: "+dest.RemoteAddr().String()+" to server:"+src.LocalAddr().String(), err)
	}
	return err
}

func startProxyTCP(listener net.Listener) error {

	log.Println("[TCP] - DNS proxy listening on TCP port 53")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("[TCP] - Error accepting connection:", err)
			conn.Close()
			continue
		}

		log.Println("[TCP] - Connecting to DNS server...")
		serverTLSConn, err := tls.Dial("tcp", dnsServerAddr+":"+dnsServerPortTLS, &config)
		if err != nil {
			log.Println("[TCP] - Error connecting to DNS server:", err)
			conn.Close()
			continue
		}

		log.Printf("[TCP] - Connected to DNS server: %s:%s", dnsServerAddr, dnsServerPortTLS)
		go copy(conn, serverTLSConn)
		go copy(serverTLSConn, conn)
		log.Printf("[TCP] - Received DNS response from %s \n", dnsServerAddr)
		// TODO: log the response from the dns server
	}
}

func main() {
	udpListener, err := net.ListenPacket("udp", ":53")
	if err != nil {
		panic(err)
	}

	tcpListener, err := net.Listen("tcp", ":53")
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()

	go startProxyTCP(tcpListener)
	go startProxyUDP(udpListener)

	defer udpListener.Close()
	defer tcpListener.Close()

	forever := make(chan bool)
	<-forever

	// TODO: Gracefully shutdown
}
