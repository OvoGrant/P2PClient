package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
)

const (
	CONFIG_FILE = "CONFIGURATION.yaml"
)

// represents a user. There's actually no real
type configStruct struct {
	DownloadLocation string `yaml:"download_location"`
}

var indexedFiles = make(map[string]bool)

var configuration configStruct

func getConfig() *configStruct {

	file, err := os.Open(CONFIG_FILE)

	if err != nil {
		createConfig()
		file, err = os.Open(CONFIG_FILE)
	}

	decoder := yaml.NewDecoder(file)

	var config configStruct

	decoder.Decode(&config)

	configuration = config

	return &config
}

func writeFileToDownloadLocation(filename string, bytes []byte) error {

	err := os.WriteFile(path.Join(configuration.DownloadLocation, filename), bytes, 0644)

	if err != nil {
		return err
	}

	return nil
}

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

	err = encoder.Encode(config)
}

// we only index new files
func getFiles() []string {

	files := make(map[string]bool)

	dir, err := os.ReadDir(configuration.DownloadLocation)

	if err != nil {
		log.Fatalf("Config error restart program")
	}

	for _, d := range dir {

		_, ok := indexedFiles[d.Name()]

		if !ok {
			fmt.Println("OK")
			indexedFiles[d.Name()] = true
			files[d.Name()] = true
		}

	}

	fileSlice := make([]string, 0, len(files))

	for k, _ := range files {
		fileSlice = append(fileSlice, k)
	}

	return fileSlice
}
