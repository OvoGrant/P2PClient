package main

import (
	"crypto/rand"
	"log"
	rand2 "math/rand"
	"net"
	"time"
)

const (
	PACKET_SIZE = 1024
)

// DelayServer contains the port that it should run on
type DelayServer struct {
	Port string
}

// NewDelayServer returns a pointer to a new delay server
func NewDelayServer(configuration *configStruct) *DelayServer {
	return &DelayServer{configuration.DelayPort}
}

// Run runs the delay server on the specified port
func (ds *DelayServer) Run() {

	//create a listener on port
	listener, err := net.Listen("tcp", ":"+ds.Port)

	//if the port is unavailable then exit program
	if err != nil {
		log.Fatal(err)
	}

	//continuously listen for connections
	for {

		conn, err := listener.Accept()

		if err != nil {
			ConnectionLogger.Println(err)
			continue
		}

		//handle the connection in a separate go routine
		go func(conn net.Conn) {

			ConnectionLogger.Println("incoming connection")

			defer conn.Close()

			//create a message of size 1024 containing random bytes
			data := makeMessage(PACKET_SIZE)

			ConnectionLogger.Println(data)

			//if the ip address is the same introduce a random delay for the sake of calculations
			if conn.RemoteAddr() == conn.LocalAddr() {
				n := rand2.Intn(4)
				time.Sleep(time.Duration(n) * time.Second)
			}

			//write the message to the connection
			conn.Write(data)
			//write a newline character
			conn.Write([]byte("\n"))

		}(conn)
	}

}

// makeMessage returns slice of random bytes of the provided size
func makeMessage(size int) []byte {
	message := make([]byte, size)
	_, err := rand.Read(message)
	if err != nil {
		panic(err)
	}
	return message
}
