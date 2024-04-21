package main

import (
	"errors"
	"io"
	"net"
)

const (
	ROOT_SERVER_TCP = "p2pcomp4911.xyz:4912"
)

// makeRequest makes to the root server for the provided filename
func makeRequest(filename string) *FileResponse {

	//dial the server for the connection
	conn, err := net.Dial("tcp", ROOT_SERVER_TCP)

	//if the error is nil create response with NETWORK_ERROR message
	if err != nil {
		return parseResponse(NETWORK_ERROR)
	}

	defer conn.Close()

	//write the name of the file with a newline character afterwards
	n, err := conn.Write([]byte(filename + "\n"))

	if err != nil {
		return parseResponse(NETWORK_ERROR)
	}

	//make buffer to store the response
	buffer := make([]byte, 1024)

	//read the response into the buffer
	n, err = conn.Read(buffer)

	if err != nil || n == 0 {
		return parseResponse("")
	}

	//return the result of parsing the response
	return parseResponse(string(buffer[:n]))

}

// initiateDownload takes in a peerResponse, file name and location to save the result to and returns an error
func initiateDownload(peer peerResponse, filename, saveLocation string) error {

	//dial the peers download port
	conn, err := net.Dial("tcp", peer.Address+":"+peer.DownloadPort)

	//if the error is nil return no connection error
	if err != nil {
		return errors.New("no connection")
	}

	//send the request over the connection
	_, err = conn.Write([]byte(filename + "\n"))

	if err != nil {
		return errors.New("download failed")
	}

	file, err := io.ReadAll(conn)

	if err != nil {
		return errors.New("download failed")
	}

	//create a new file util to save the file to disk
	fileUtil := NewFileUtil(saveLocation)

	//write the file to location
	err = fileUtil.writeFileToDownloadLocation(filename, file)

	if err != nil {
		return errors.New("error writing file to disk")
	}

	return nil
}
