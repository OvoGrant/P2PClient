package main

import (
	"log"
	"os"
	"path"
)

// type fileUtil contains the files that have been indexed and their location in the file system
type fileUtil struct {
	indexedFiles map[string]bool
	location     string
}

// NewFileUtil returns an instance to a new fileUtil
func NewFileUtil(DownloadLocation string) *fileUtil {
	return &fileUtil{make(map[string]bool), DownloadLocation}
}

// writeFileToDownloadLocation writes a slice of bytes to a file with the provided name
func (fu *fileUtil) writeFileToDownloadLocation(filename string, bytes []byte) error {

	err := os.WriteFile(path.Join(fu.location, filename), bytes, 0644)

	if err != nil {
		return err
	}

	return nil
}

// getFiles returns a slice of filenames that have not been indexed yet
func (fu *fileUtil) getFiles() []string {

	files := make(map[string]bool)

	dir, err := os.ReadDir(fu.location)

	if err != nil {
		log.Fatalf("Config error restart program")
	}

	for _, d := range dir {

		_, ok := fu.indexedFiles[d.Name()]

		if !ok {
			fu.indexedFiles[d.Name()] = true
			files[d.Name()] = true
		}

	}

	fileSlice := make([]string, 0, len(files))

	for k, _ := range files {
		fileSlice = append(fileSlice, k)
	}

	return fileSlice
}
