package main

import (
	"encoding/json"
	"net"
	"strings"
	"time"
)

const (
	// Note that this can be changed to the localhost if the indexing server is running locally
	ROOT_SERVER_UDP = "p2pcomp4911.xyz:4911"
	UDP             = "udp4"
)

// IndexingClient reads the files in the download location and notifies the central server about them
type IndexingClient struct {
	FileUtil *fileUtil
	Peer     Peer
}

// NewIndexingClient returns a pointer to an IndexingClient
func NewIndexingClient(configuration *configStruct) *IndexingClient {
	var client IndexingClient

	client.Peer = Peer{configuration.DownloadPort, configuration.DelayPort}

	client.FileUtil = NewFileUtil(configuration.DownloadLocation)

	return &client
}

// Run runs the IndexingClient
func (ic *IndexingClient) Run() {

	address, err := net.ResolveUDPAddr(UDP, ROOT_SERVER_UDP)

	conn, err := net.DialUDP(UDP, nil, address)

	if err != nil {
		ErrorLogger.Println("Error could not connect to internet")
	}

	for {

		time.Sleep(1 * time.Second)

		files := ic.FileUtil.getFiles()

		bytes, err := json.Marshal(ic.Peer)

		if len(files) > 0 {
			bytes = append(bytes, []byte("\n"+strings.Join(files, "\n"))...)
		}

		_, err = conn.Write(bytes)

		conn.SetReadDeadline(time.Now().Add(4 * time.Second))

		response := make([]byte, 1024)

		n, _, err := conn.ReadFrom(response)

		if err != nil {
			ErrorLogger.Fatal(err)
			continue
		}

		resp := string(response[0:n])

		if resp == ACK {
			time.Sleep(5 * time.Second)
		}
	}
}
