package core

import (
	"encoding/json"
	"net/http"

	"github.com/etowett/go-api-sim/utils"
)

type AFTResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func ATPage(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	destinaton := r.FormValue("to")
	message := r.FormValue("message")
	from := r.FormValue("from")
	key := r.FormValue("key")

	req := map[string]string{
		"username": username, "to": destinaton, "message": message,
		"from": from, "key": key,
	}

	utils.Logger.Println("Request: ", req)

	json.NewEncoder(w).Encode(AFTResponse{
		Status: "success", Message: "Request Received",
	})
}
