package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

const (
	CONFIG_FILE = "CONFIGURATION.yaml"
)

// represents a user.
type configStruct struct {
	DownloadLocation string `yaml:"download_location"`
	UploadLocation   string `yaml:"upload_location""`
	DownloadPort     string `yaml:"download_port"`
	DelayPort        string `yaml:"delay_port"`
}

// IsValid ensures that configStruct contains the correct information
func (cs *configStruct) IsValid() bool {
	dp, err := strconv.Atoi(cs.DelayPort)

	if err != nil {
		return false
	}

	dl, err := strconv.Atoi(cs.DownloadPort)

	if err != nil {
		return false
	}

	if dp < 0 || dl < 0 {
		return false
	}

	if len(cs.DownloadLocation) == 0 {
		return false
	}

	dir, err := os.Stat(cs.DownloadLocation)

	if err != nil || !dir.IsDir() {
		return false
	}

	return true
}

var configuration configStruct

// getConfig returns a configStruct by parsing the information inside of CONFIGURATION.YML
func getConfig() *configStruct {

	file, err := os.Open(CONFIG_FILE)

	if err != nil {
		createConfig()
		file, err = os.Open(CONFIG_FILE)
	}

	decoder := yaml.NewDecoder(file)

	var config configStruct

	decoder.Decode(&config)

	for !config.IsValid() {

		createConfig()

		file, err = os.Open(CONFIG_FILE)

		decoder.Decode(&config)

		configuration = config

	}

	return &config
}

// createConfig prompts the user for input in order to create a config file
func createConfig() {

	file, err := os.OpenFile(CONFIG_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		log.Fatalf("Error creating file %v", err)
	}

	defer file.Close()

	encoder := yaml.NewEncoder(file)

	var config configStruct

	fmt.Println("Enter the download location you wish to use to store your files")

	var dirName string

	fmt.Scanln(&dirName)

	dir, err := os.Stat(dirName)

	for err != nil || !dir.IsDir() {

		fmt.Println("Enter the download location you wish to use to store your files")

		fmt.Scanln(&dirName)

		dir, err = os.Stat(dirName)

	}

	config.DownloadLocation = dirName

	fmt.Println("Please enter the port you want to handle downloads on")

	downloadPort, err := getInteger()

	for err != nil {
		fmt.Println("Please enter the port you want to handle downloads on")
		downloadPort, err = getInteger()
	}

	config.DownloadPort = strconv.Itoa(downloadPort)

	fmt.Println("Please enter the port you to use to handle delay calculations")
	delayPort, err := getInteger()

	for err != nil {
		fmt.Println("Please enter the port you to use to handle delay calculations")
		delayPort, err = getInteger()
	}

	config.DelayPort = strconv.Itoa(delayPort)

	err = encoder.Encode(config)
}

func getInteger() (int, error) {
	var i int
	_, err := fmt.Scanf("%d", &i)

	if err != nil {
		return i, err
	}

	return i, nil
}
