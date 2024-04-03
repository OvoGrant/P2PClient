package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

const (
	DOWNLOAD_PORT   = ":4200"
	TEST_PORT       = ":4100"
	ROOT_SERVER_UDP = "p2pcomp4911.xyz:4911"
	ROOT_SERVER_TCP = "p2pcomp4911.xyz:4912"
)

func makeRequest(filename string) *FileResponse {

	conn, err := net.Dial("tcp", ROOT_SERVER_TCP)

	if err != nil {
		return parseResponse(NETWORK_ERROR)
	}

	defer conn.Close()

	n, err := conn.Write([]byte(filename + "\n"))

	if err != nil {
		return parseResponse(NETWORK_ERROR)
	}

	buffer := make([]byte, 1024)

	n, err = conn.Read(buffer)

	if err != nil || n == 0 {
		return parseResponse(" ")
	}

	return parseResponse(string(buffer[:n]))

}

func initiateDownload(peerAddress string, filename string) error {

	fmt.Println(peerAddress + DOWNLOAD_PORT)
	fmt.Println(len(peerAddress + DOWNLOAD_PORT))

	conn, err := net.Dial("tcp", peerAddress+DOWNLOAD_PORT)

	if err != nil {
		fmt.Println(err)
		return errors.New("no connection")
	}

	//send the request over the connection
	conn.Write([]byte(filename + "\n"))

	file, err := io.ReadAll(conn)

	if err != nil {
		return errors.New("download failed")
	}

	err = writeFileToDownloadLocation(filename, file)

	if err != nil {
		return errors.New("error writing file to disk")
	}

	return nil
}

// this should handle the download
func downloadServer() {

	listener, err := net.Listen("tcp", "127.0.0.1"+DOWNLOAD_PORT)

	ConnectionLogger.Println("Starting download server...")

	defer listener.Close()

	if err != nil {
		ConnectionLogger.Println("Error: \n", err)
	}

	for {

		conn, err := listener.Accept()

		if err != nil {
			ConnectionLogger.Println(err)
			continue
		}

		go handleDownload(conn)
	}
}

// read file and download it
func handleDownload(conn net.Conn) {

	ConnectionLogger.Printf("Incoming connection from %s\n", conn.LocalAddr())

	defer conn.Close()

	//a download request should consist of one line that is the name of the file
	reader := bufio.NewScanner(conn)

	if !reader.Scan() {
		fmt.Println("Error: ", reader.Err())
		conn.Write([]byte("NACK\nFormat Error"))
		return
	}

	fmt.Println("debug")

	filename := reader.Text()

	fmt.Println(filename)

	ConnectionLogger.Println(filename)

	file, err := os.ReadFile(path.Join(configuration.DownloadLocation, filename))

	fmt.Println(file)

	if err != nil {
		fmt.Println(err)
		conn.Write([]byte("Nack\n" + err.Error()))
	}

	nn, err := conn.Write(file)

	fmt.Println(nn)

}

func indexingClient() {

	address, err := net.ResolveUDPAddr("udp4", ROOT_SERVER_UDP)

	conn, err := net.DialUDP("udp4", nil, address)

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error could not connect to internet")
	}

	//the indexing client works theoretically
	for {

		time.Sleep(1 * time.Second)

		files := getFiles()

		_, err = conn.Write([]byte(strings.Join(files, "\n")))

		response := make([]byte, 1024)

		n, _, err := conn.ReadFrom(response)

		if err != nil {
			log.Fatal(err)
		}

		resp := string(response[0:n])

		if resp == ACK {
			time.Sleep(5 * time.Second)
		}
	}
}
