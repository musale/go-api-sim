package main

import (
	"log"
	"net/http"
	"os"

	"github.com/etowett/go-api-sim/src"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env Error ", err)
	}

	f, err := os.OpenFile(
		os.Getenv("LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600,
	)
	if err != nil {
		log.Fatal("Log file error: ", err)
	}
	defer f.Close()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(f)

	spinATWorkers()
	spinRMWorkers()

	// Route set up
	http.HandleFunc("/aft", src.ATPage)
	http.HandleFunc("/routesms", src.RMPage)
	http.HandleFunc("/saf", src.SafPage)
	http.HandleFunc("/saf-dlr", src.SafDLRPage)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func spinATWorkers() {
	for i := 0; i < 250; i++ {
		go func() {
			for req := range src.ATReqChan {
				result := src.ProcessATReq(&req)
				src.ATResChan <- result
			}
		}()
	}
}

func spinRMWorkers() {
	for i := 0; i < 250; i++ {
		go func() {
			for req := range src.RMReqChan {
				result := src.ProcessRMReq(&req)
				src.RMResChan <- result
			}
		}()
	}
}
