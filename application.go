package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// constants used for prompting the user for information
const (
	WELCOME_MESSAGE = "Hello Welcome to my p2p file-sharing application"
	FILE_PROMPT     = "Hi type in the name of a file that you wish to get press q to quit"
	QUIT            = "q"
)

// the process interface contains only one method Run which will be invoked in a go routine
type process interface {
	Run()
}

// application struct contains a configuration struct and a list of background process to run
type application struct {
	configuration     *configStruct
	backgroundProcess []process
}

// start gets the configuration adds the necessary processes to the slice of processes and runs them in the background
// then it runs the main foreground process
func (ap *application) start() {

	initLoggers()

	//get configuration
	ap.configuration = getConfig()

	//add the background services that we want to run
	ap.backgroundProcess = append(ap.backgroundProcess, NewIndexingClient(ap.configuration))
	ap.backgroundProcess = append(ap.backgroundProcess, NewDelayServer(ap.configuration))
	ap.backgroundProcess = append(ap.backgroundProcess, NewDownloadServer(ap.configuration))

	//for each process run it in the background
	for _, process := range ap.backgroundProcess {
		go process.Run()
	}

	//run the main foreground application
	ap.RunApplication()
}

// RunApplication continuously prompts the user for input and completes actions based on the input
func (ap *application) RunApplication() {

	//print the welcome message
	fmt.Println(WELCOME_MESSAGE)

	//create a scanner
	scanner := bufio.NewScanner(os.Stdin)

	var filename string

	//run while the user has not requested to quit the program
	for filename != QUIT {

		//print the file prompt
		fmt.Println(FILE_PROMPT)

		scanner.Scan()
		filename = scanner.Text()

		//make the request for the file
		response := makeRequest(filename)

		if response.OK() {

			//if the file response is ok then print the list of peers and their associated information
			fmt.Println("The following peers have the file")

			dc := newPeerStatsCalculator(response.Peers)

			//calculate the delay and latency for each peer
			delays := dc.calculateDelays()

			//print out the peer information in a formatted table
			printPeerTable(delays)

			fmt.Println("Select which peer you would like to download the file from")

			var line string

			scanner.Scan()
			line = scanner.Text()

			//get the peer that the user wishes to get the file from
			peerNumber, err := strconv.Atoi(line)

			for err != nil && !(peerNumber < 0 || peerNumber >= len(response.Peers)) {
				fmt.Println("Select which peer you would like to connect with (Enter a valid index)")

				scanner.Scan()
				line = scanner.Text()

				peerNumber, err = strconv.Atoi(line)
			}

			//initiate download
			err = initiateDownload(response.GetPeer(peerNumber), filename, ap.configuration.DownloadLocation)

			//if there is an error during the download print it and continue the process
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			//print download success
			fmt.Println("Download was a success")

		} else {

			fmt.Println(response.Error)

		}

	}

}
