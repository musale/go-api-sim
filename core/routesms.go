package core

import (
	"fmt"
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

	w.Header().Set("Server", "G-Starfish")
	w.WriteHeader(200)

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

	if validateUser(username, password) == false {
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
		valid, num := phone.IsValid(number)
		var x string
		if valid == false {
			x = "1706|" + number
		} else {
			x = "1701|" + num[1:] + "|" + utils.GetUUID()
		}
		data = append(data, x)
	}

	utils.Logger.Println("RMS Message: ", message)
	utils.Logger.Println("RMS Recipients: ", len(strings.Split(destinaton, ",")))

	fmt.Fprintf(w, strings.Join(data, ","))
	return
}
