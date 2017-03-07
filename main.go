package main

import (
	"log"
	"net/http"
	"os"

	"github.com/etowett/go-api-sim/core"
	"github.com/etowett/go-api-sim/utils"
	"github.com/joho/godotenv"
)

var err error

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file ", err)
		return
	}

	logFile, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		log.Println("Log file error: ", err)
	}

	defer logFile.Close()

	utils.Logger = log.New(logFile, "", log.Lshortfile|log.Ldate|log.Ltime)

	// Route set up
	http.HandleFunc("/aft", core.ATPage)
	http.HandleFunc("/routesms", core.RMPage)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
