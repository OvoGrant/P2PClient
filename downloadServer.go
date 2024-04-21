package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"path"
)

// DownloadServer contains a location to search for files to and a port to handle requests on
type DownloadServer struct {
	DownloadLocation string
	Port             string
}

// NewDownloadServer takes in a configStruct returns a pointer to a DownloadServer
func NewDownloadServer(configuration *configStruct) *DownloadServer {
	return &DownloadServer{configuration.DownloadLocation, configuration.DownloadPort}
}

// Run runs the server logic
func (ds *DownloadServer) Run() {

	//create a listener
	listener, err := net.Listen("tcp", ":"+ds.Port)

	ConnectionLogger.Println("Starting download server...")

	defer listener.Close()

	//if there is an error exit the program
	if err != nil {
		log.Fatal("Error: \n", err)
	}

	//constantly listen for connections
	for {

		//accept the connection
		conn, err := listener.Accept()

		ConnectionLogger.Println("Incoming connection from " + conn.RemoteAddr().String())

		if err != nil {
			ConnectionLogger.Println(err)
			continue
		}

		go ds.handleDownload(conn)
	}
}

// read file and download it
func (ds *DownloadServer) handleDownload(conn net.Conn) {

	ConnectionLogger.Printf("Incoming connection from %s\n", conn.LocalAddr())

	defer conn.Close()

	//a download request should consist of one line that is the name of the file
	reader := bufio.NewScanner(conn)

	if !reader.Scan() {
		ErrorLogger.Println("Error: ", reader.Err())
		conn.Write([]byte("NACK\nFormat Error"))
		return
	}

	filename := reader.Text()

	ConnectionLogger.Println(filename)

	file, err := os.ReadFile(path.Join(ds.DownloadLocation, filename))

	if err != nil {
		ErrorLogger.Println(err)
		conn.Write([]byte("Nack\n" + err.Error()))
	}

	_, err = conn.Write(file)
}
