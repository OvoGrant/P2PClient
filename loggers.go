package main

import (
	"log"
	"os"
)

var (
	ConnectionLogger *log.Logger
	ErrorLogger      *log.Logger
)

func initLoggers() {

	file, err := os.Create("logs.txt")

	if err != nil {
		log.Fatalf(err.Error())
	}

	ConnectionLogger = log.New(file, "CONNECTION: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}
