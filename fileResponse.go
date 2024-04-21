package main

import (
	"encoding/json"
	"strings"
)

const (
	ACK           = "ACK"
	NACK          = "NACK"
	NETWORK_ERROR = "NETWORK_ERROR"
)

// FileResponse represents the message returned after a client requests a file from the indexing server
type FileResponse struct {
	Status string
	Error  string
	Peers  []peerResponse
}

// parseResponse parses a string and returns a fileResponse struct
func parseResponse(data string) *FileResponse {
	var response FileResponse

	body := strings.Split(data, "\n")

	if body[0] == ACK {

		response.Status = ACK

		response.Peers = make([]peerResponse, 0)

		for i := 1; i < len(body); i++ {
			peer := peerResponse{}

			json.Unmarshal([]byte(body[i]), &peer)

			response.Peers = append(response.Peers, peer)
		}

		return &response

	} else if body[0] == NACK {

		response.Status = NACK
		response.Error = body[1]
		return &response

	} else {
		response.Status = NETWORK_ERROR
		response.Error = NETWORK_ERROR
	}

	return &response
}

// OK returns whether the result of the request was good
func (fr *FileResponse) OK() bool {
	return fr.Status == ACK
}

// GetPeer returns the peer at a particular position
func (fr *FileResponse) GetPeer(position int) peerResponse {
	return fr.Peers[position]
}
