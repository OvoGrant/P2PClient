package main

import (
	"io"
	"net"
	"sort"
	"time"
)

const (
	PACKET_SIZE = 1024
)

type peerStatsCalculator struct {
	peers []string
}

type peerStats interface {
	calculateThroughput() []peerDelay
}

func newPeerStatsCalculator(peers []string) *peerStatsCalculator {
	var dc peerStatsCalculator

	for _, peer := range peers {
		dc.peers = append(dc.peers, peer)
	}

	return &dc
}

func (dc *peerStatsCalculator) calculateDelays() []peerDelay {

	results := make(chan peerDelay)

	for _, host := range dc.peers {
		go determineDelay(host, TEST_PORT, results)
	}

	delayResult := make([]peerDelay, 0)

	for i := 0; i < len(dc.peers); i++ {
		result := <-results
		delayResult = append(delayResult, result)
	}

	sort.Slice(delayResult, func(i, j int) bool {
		return delayResult[i].rate < delayResult[j].rate
	})

	return delayResult

}

func determineDelay(host string, port string, delay chan peerDelay) {

	var stats peerDelay
	stats.addr = host

	latencyStart := time.Now()

	conn, err := net.DialTimeout("tcp", host+port, 5*time.Second)

	if err != nil {
		stats.rate = PACKET_SIZE
		delay <- stats
		return
	}

	defer conn.Close()

	stats.latency = time.Now().Sub(latencyStart).Milliseconds()

	startTime := time.Now()

	_, err = io.ReadAll(conn)

	if err != nil {
		delay <- stats
		return
	}

	endTime := time.Now()

	stats.rate = float32(PACKET_SIZE) / float32(endTime.Sub(startTime).Seconds())

	delay <- stats
}

type peerDelay struct {
	addr    string
	rate    float32
	latency int64
}
