package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	WELCOME_MESSAGE = "Hello Welcome to my p2p file-sharing application"
	FILE_PROMPT     = "Hi type in the name of a file that you wish to get press q to quit"
	QUIT            = "q"
)

func RunApplication() {

	fmt.Println(WELCOME_MESSAGE)

	scanner := bufio.NewScanner(os.Stdin)

	var filename string

	for filename != QUIT {

		fmt.Println(FILE_PROMPT)

		scanner.Scan()
		filename = scanner.Text()

		response := makeRequest(filename)

		if response.OK() {

			fmt.Println("The following peers have the file")

			dc := newPeerStatsCalculator(response.Peers)

			delays := dc.calculateDelays()

			printPeerTable(delays)

			fmt.Println("Select which peer you would like to download the file from")

			var line string

			scanner.Scan()
			line = scanner.Text()

			peerNumber, err := strconv.Atoi(line)

			if err != nil {
			}

			err = initiateDownload(response.GetPeer(peerNumber), filename)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Println("Download was a success")

		} else {

			fmt.Println("No users had the file you requested")

		}

	}

}
