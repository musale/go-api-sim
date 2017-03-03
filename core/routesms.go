package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/etowett/go-api-sim/phone"
)

type RMResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func RMPage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Server", "G-Starfish")
	w.WriteHeader(200)

	username := r.FormValue("username")
	password := r.FormValue("password")
	message := r.FormValue("message")
	from := r.FormValue("source")
	destinaton := r.FormValue("destination")
	dlr := r.FormValue("dlr")
	typ := r.FormValue("type")

	if message == "" || len(message) == 0 {
		log.Println("No message")
		fmt.Fprintf(w, "1702\n")
		return
	}

	if from == "" || len(from) == 0 {
		log.Println("No from")
		fmt.Fprintf(w, "1702\n")
		return
	}

	// if password == "" || len(password) == 0 || username == "" ||
	// 	len(username) == 0 || from == "" || len(from) == 0 ||
	// 	destinaton == "" || len(destinaton) == 0 || message == "" ||
	// 	len(message) == 0 || dlr == "" || len(dlr) == 0 ||
	// 	typ == "" || len(typ) == 0 {
	// 	fmt.Fprintf(w, "1702\n")
	// 	return
	// }

	if validateUser(username, password) == false {
		fmt.Fprintf(w, "1703.\n")
		return
	}

	if typ != "0" {
		fmt.Fprintf(w, "1704.\n")
		return
	}

	if dlr != "0" || dlr != "1" {
		fmt.Fprintf(w, "1708.\n")
		return
	}

	for _, number := range strings.Split(destinaton, ",") {
		valid, num := phone.IsValid(number)
		if valid == false {
			log.Println("valid false", valid)
		}
		log.Println("Number ", num)
	}

	json.NewEncoder(w).Encode(AFTResponse{
		Status: "success", Message: "Request Received",
	})
}
