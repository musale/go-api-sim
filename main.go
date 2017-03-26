package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/etowett/go-api-sim/core"
	"github.com/etowett/go-api-sim/utils"
	"github.com/joho/godotenv"
)

var (
	err     error
	logpath = flag.String("logpath", "/var/log/go-api/api.log", "Log Path")
)

func main() {

	err = godotenv.Load()
	if err != nil {
		log.Fatal(".env Error ", err)
	}

	flag.Parse()
	utils.NewLog(*logpath)
	utils.Log.Println("hello there")

	// f, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	log.Fatal("log file error: ", err)
	// }
	// defer f.Close()

	// myFile := log.New(f,
	// 	"PREFIX: ",
	// 	log.Ldate|log.Ltime|log.Lshortfile)

	// log.SetOutput(myFile)

	// Route set up
	http.HandleFunc("/aft", core.ATPage)
	http.HandleFunc("/routesms", core.RMPage)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
