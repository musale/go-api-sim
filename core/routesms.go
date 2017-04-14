package core

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/etowett/go-api-sim/phone"
	"github.com/etowett/go-api-sim/utils"
)

type RMResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func RMPage(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	message := r.FormValue("message")
	from := r.FormValue("source")
	destinaton := r.FormValue("destination")
	dlr := r.FormValue("dlr")
	typ := r.FormValue("type")

	if password == "" || len(password) == 0 || username == "" ||
		len(username) == 0 || from == "" || len(from) == 0 ||
		destinaton == "" || len(destinaton) == 0 || message == "" ||
		len(message) == 0 || dlr == "" || len(dlr) == 0 ||
		typ == "" || len(typ) == 0 {
		fmt.Fprintf(w, "1702\n")
		return
	}

	if !validateUser(username, password) {
		fmt.Fprintf(w, "1703\n")
		return
	}

	if typ != "0" {
		fmt.Fprintf(w, "1704\n")
		return
	}

	if dlr != "0" {
		if dlr != "1" {
			fmt.Fprintf(w, "1708\n")
			return
		}
	}

	var data []string
	for _, number := range strings.Split(destinaton, ",") {
		num, err := phone.CheckValid(number)
		var x string
		if err != nil {
			x = "1706|" + number
		} else {
			x = "1701|" + num[1:] + "|" + utils.GetUUID()
		}
		data = append(data, x)
	}

	log.Println("RMS Message: ", message)
	log.Println("RMS Recipients: ", len(strings.Split(destinaton, ",")))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strings.Join(data, ","))
	return
}
