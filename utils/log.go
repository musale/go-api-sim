package utils

import (
	"log"
	"os"
)

var (
	// Log Common Logger var
	Log *log.Logger
)

// NewLog Set up log
func NewLog(logpath string) {
	println("LogFile: " + logpath)
	file, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}
