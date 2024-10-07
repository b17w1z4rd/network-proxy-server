package main

import (
	"io"
	"log"
	"net"
	"os"
)

func handleConnection(src net.Conn, targetAddr string) {
	defer src.Close()

	dst, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Error connecting to target %s: %v", targetAddr, err)
		return
	}
	defer dst.Close()

	go io.Copy(dst, src) // Copy data from the source to the destination
	io.Copy(src, dst)    // Copy data from the destination to the source
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <listen_addr> <target_addr>", os.Args[0])
	}

	listenAddr := os.Args[1]
	targetAddr := os.Args[2]

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening on %s and forwarding to %s", listenAddr, targetAddr)

	for {
		src, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from %s", src.RemoteAddr())
		go handleConnection(src, targetAddr) // Handle connection concurrently
	}
}
