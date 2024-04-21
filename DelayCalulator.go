package main

import (
	"io"
	"net"
	"sort"
	"time"
)

// peerDelay contains the address, download rate and latency associated with a particular peer
type peerDelay struct {
	addr    string
	rate    float32
	latency int64
}

// peerStatsCalculator contains a slice of peers
type peerStatsCalculator struct {
	peers []peerResponse
}

// A peerStats calculator is a generic interface that calculates the rate and latency for a connection
type peerStats interface {
	calculateThroughput() []peerDelay
}

// newPeerStatsCalculator returns a reference to a new peerStatsCalculator
func newPeerStatsCalculator(peers []peerResponse) *peerStatsCalculator {
	var dc peerStatsCalculator

	for _, peer := range peers {
		dc.peers = append(dc.peers, peer)
	}

	return &dc
}

// calculateDelays calculates the delay and latency for each peer in the peerStatsCalculator
func (dc *peerStatsCalculator) calculateDelays() []peerDelay {

	//create channel to read results
	results := make(chan peerDelay)

	//for each peer in the peer slice calculate stats in separate go routine
	for _, host := range dc.peers {
		go determineDelay(host.getAddress(), host.getDelayPort(), results)
	}

	//make slice to store results
	delayResult := make([]peerDelay, 0)

	//for the number of peers in the slice wait until one writes to channel and then append it to the final result
	for i := 0; i < len(dc.peers); i++ {
		result := <-results
		delayResult = append(delayResult, result)
	}

	//sort the slice with the greatest rate coming first
	sort.Slice(delayResult, func(i, j int) bool {
		return delayResult[i].rate > delayResult[j].rate
	})

	return delayResult

}

// determineDelay calculates the latency and delay for a particular host
func determineDelay(host string, port string, delay chan peerDelay) {

	var stats peerDelay
	stats.addr = host

	//start the latency timer
	latencyStart := time.Now()

	//dial the host
	conn, err := net.DialTimeout("tcp", host+":"+port, 5*time.Second)

	//if there is an error establishing the connection then set default rate and return
	if err != nil {
		stats.rate = 0
		delay <- stats
		return
	}

	defer conn.Close()

	//set the value of latency
	stats.latency = time.Now().Sub(latencyStart).Milliseconds()

	//star the download timer
	startTime := time.Now()

	//read the bytes from the connection
	_, err = io.ReadAll(conn)

	//if there is an error reading from the connection then return the struct
	if err != nil {
		delay <- stats
		return
	}

	//get end time
	endTime := time.Now()

	//set the rate equal to the packet size divided by the time to read it
	bytesPerSecond := float32(PACKET_SIZE) / float32(endTime.Sub(startTime).Seconds())

	//divide by 1000000 to get Megabits per second
	stats.rate = bytesPerSecond / 1000000

	delay <- stats
}
