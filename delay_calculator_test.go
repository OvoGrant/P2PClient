package main

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

var dummy_bytes = make([]byte, 1024)

func TestPeerStatCalculator(t *testing.T) {

	port := ":5944"

	go testServers(port)

	result := make(chan peerDelay)

	go determineDelay("", port, result)

	res := <-result

	printPeerTable([]peerDelay{res})

}

func testServers(port string) {

	listener, err := net.Listen("tcp", "127.0.0.1"+port)

	defer listener.Close()

	if err != nil {
		return
	}

	success := false

	for {

		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		go func(net.Conn) {

			defer conn.Close()
			defer listener.Close()

			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

			fmt.Println("writing...")
			if _, err := conn.Write(dummy_bytes); err == nil {
				success = true
			}

		}(conn)
		if success {
			break
		}

	}
}
