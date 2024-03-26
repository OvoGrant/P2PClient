package main

import (
	"fmt"
	"strings"
)

const (
	ACK           = "ACK"
	NACK          = "NACK"
	NETWORK_ERROR = "NETWORK_ERROR"
)

type FileResponse struct {
	Status string
	Error  string
	Peers  []string
}

func parseResponse(data string) *FileResponse {
	var response FileResponse

	body := strings.Split(data, "\n")

	fmt.Println(body)

	if body[0] == ACK {

		response.Status = ACK

		response.Peers = make([]string, 0)

		for i := 1; i < len(body); i++ {
			response.Peers = append(response.Peers, body[i])
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

func (fr *FileResponse) OK() bool {
	return fr.Status == ACK
}

func (fr *FileResponse) GetPeer(position int) string {
	return fr.Peers[position]
}
