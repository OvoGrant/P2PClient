package main

// peerResponse represent the response a client receives when requesting a list of peers
type peerResponse struct {
	Address string `yaml:"address" json:"address"`
	Peer
}

// Peer represents the information a client sends to the central server
type Peer struct {
	DownloadPort string `yaml:"download_port" json:"download_port"`
	DelayPort    string `yaml:"delay_port" json:"delay_port"`
}

// getAddress returns the ip address of a peer
func (pr *peerResponse) getAddress() string {
	return pr.Address
}

// getDownloadPort returns the port that the peer accepts downloads on
func (pr *Peer) getDownloadPort() string {
	return pr.DownloadPort
}

// getDelayPort returns the port on which the peer accepts requests for delay calculation
func (pr *Peer) getDelayPort() string {
	return pr.DelayPort
}
