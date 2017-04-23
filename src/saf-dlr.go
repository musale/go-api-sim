package src

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/etowett/go-api-sim/utils"
)

// SafDLRPage endpoint to rm request
func SafDLRPage(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ioutil readall: ", err)
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.Write([]byte("read all error"))
		return
	}

	var req SMSEnvelope
	if err := xml.Unmarshal(body, &req); err != nil {
		log.Println("Xml unmarshal: ", err)
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.Write([]byte("xml unmarshal err"))
		return
	}

	reqBody := req.SMSReqBody.RequestBody
	log.Println(fmt.Sprintf(
		"Request:: SenderID: %s, Phone: %s, Message: %s", reqBody.SenderID,
		reqBody.Number, reqBody.Message,
	))

	senderIDs := []string{"FOCUSMOBILE", "Eutychus", "SMSLEOPARD"}

	if !utils.InArray(reqBody.SenderID, senderIDs) {
		faultCode := "SVC0002"
		faultString := "SenderName or senderAddress is unknown!"
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.Write([]byte(fmt.Sprintf(
			senderIDResponse, faultCode, faultString, faultCode, faultString,
			reqBody.SenderID,
		)))
		return
	}
	msgID := "100001200501170419072620015931"
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
	w.Write([]byte(fmt.Sprintf(successResponse, msgID)))
	return
}
