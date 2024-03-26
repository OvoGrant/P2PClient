package main

import (
	"log"
	"os"
)

var (
	ConnectionLogger *log.Logger
	DeletionLogger   *log.Logger
)

func initLoggers() {

	file, err := os.Open("logs.txt")

	if err != nil {
		log.Fatalf(err.Error())
	}

	ConnectionLogger = log.New(file, "CONNECTION: ", log.Ldate|log.Ltime|log.Lshortfile)
	DeletionLogger = log.New(file, "DELETION: ", log.Ldate|log.Ltime|log.Lshortfile)
}
