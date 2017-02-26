package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/etowett/go-api-sim/phone"
	"github.com/etowett/go-api-sim/utils"
)

type Recipient struct {
	Number    string  `json:"number"`
	Cost      float64 `json:"cost"`
	Status    string  `json:"status"`
	MessageId string  `json:"messageId"`
}

type MessageData struct {
	Message    string      `json:"Message"`
	Recipients []Recipient `json:"Recipients"`
}

type FinalResponse struct {
	SMSMessageData MessageData `json:"SMSMessageData"`
}

type AFTResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func ATPage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Server", "G-Starfish")
	w.WriteHeader(200)

	username := r.FormValue("username")
	destinaton := r.FormValue("to")
	message := r.FormValue("message")
	from := r.FormValue("from")
	key := r.Header.Get("apikey")

	if key == "" || len(key) == 0 {
		fmt.Fprintf(w, "Request is missing required HTTP header apikey.\n")
		return
	}

	if username == "" || len(username) == 0 {
		fmt.Fprintf(w, "Must have username in your request.\n")
		return
	}

	if validateUser(username, key) == false {
		fmt.Fprintf(w, "The supplied authentication is invalid.\n")
		return
	}

	if from == "" || len(from) == 0 {
		fmt.Fprintf(w, "Must have from in your request.\n")
		return
	}

	if destinaton == "" || len(destinaton) == 0 {
		fmt.Fprintf(w, "Must have to in your request.\n")
		return
	}

	if message == "" || len(message) == 0 {
		fmt.Fprintf(w, "Must have message in your request.\n")
		return
	}

	var recipients []Recipient

	for _, number := range strings.Split(destinaton, ",") {
		var rec Recipient
		valid, num := phone.IsValid(number)
		cost := 0.0
		status := "Failed"
		messageID := "None"
		if valid == false {
			status = "Invalid Phone Number"
		} else if inBlacklist(number) == true {
			status = "User In Blacklist"
		} else {
			number = num
			status = "Success"
			cost = getMesageCost(message, num)
			messageID = utils.GetMD5Hash(time.Now().String() + number)
		}
		rec.Status = status
		rec.Cost = cost
		rec.Number = number
		rec.MessageId = messageID
		recipients = append(recipients, rec)
	}

	msg := "Sent to 1/7 Total Cost: KES 500"

	ret := FinalResponse{
		SMSMessageData: MessageData{
			Message: msg, Recipients: recipients,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}

func validateUser(username, key string) bool {
	return true
}

func inBlacklist(number string) bool {
	values := []bool{true, false}
	return values[rand.Intn(len(values))]
}

func getMesageCost(message string, number string) float64 {
	return 1.0
}
