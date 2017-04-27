package src

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/etowett/go-api-sim/phone"
	"github.com/etowett/go-api-sim/utils"
)

// Recipient single destination data
type Recipient struct {
	Number    string `json:"number"`
	Cost      string `json:"cost"`
	Status    string `json:"status"`
	MessageID string `json:"messageId"`
}

// MessageData
type MessageData struct {
	Message    string      `json:"Message"`
	Recipients []Recipient `json:"Recipients"`
}

// FinalResponse format for final response
type FinalResponse struct {
	SMSMessageData MessageData `json:"SMSMessageData"`
}

// AFTResponse payload for response
type AFTResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// ATPage handler for AT request
func ATPage(w http.ResponseWriter, r *http.Request) {

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

	go func() {
		utils.ATChan <- ATRequest{
			Destinaton: Number, Message: message,
		}
	}()

	smsData := <-utils.ATResult

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FinalResponse{smsData})
	return
}

func ProcessATReq(req *ATRequest) MessageData {
	var recipients []Recipient
	validCount := 0
	totalCost := 0.0
	for _, number := range strings.Split(req.Destinaton, ",") {
		var rec Recipient
		var cost = 0.0
		var status = "Failed"
		var messageID = "None"

		num, err := phone.CheckValid(number)
		if err != nil {
			status = "Invalid Phone Number"
		} else if !inBlacklist(number) {
			status = "User In BlackList"
		} else {
			number = num
			status = "Success"
			cost, _ = getMesageCost(message, num)
			messageID = utils.GetMD5Hash(time.Now().String() + number)
			validCount++
			totalCost += cost
		}
		rec.Status = status
		rec.Cost = fmt.Sprintf("%.2f", cost)
		rec.Number = number
		rec.MessageID = messageID
		recipients = append(recipients, rec)
	}
	log.Println("AFT Message: ", req.Message)
	log.Println("AFT Recipients: ", len(strings.Split(req.Destinaton, ",")))

	msg := fmt.Sprintf(
		"Sent to %v/%v Total Cost: KES %v", validCount, len(recipients), totalCost)
	return MessageData{
		Message: msg, Recipients: recipients,
	}
}

func validateUser(username, key string) bool {
	return true
}

func inBlacklist(number string) bool {
	values := []bool{true, false}
	return values[rand.Intn(len(values))]
}

func getMesageCost(message string, number string) (float64, error) {
	return 1.0, nil
}
