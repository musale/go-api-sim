package src

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	// "math/rand"
	"net/http"
)

var senderIDResponse = `<?xml version="1.0" encoding="UTF-8"?><soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><soapenv:Body><soapenv:Fault><faultcode>%s</faultcode><faultstring>%s</faultstring><detail><ns1:ServiceException xmlns:ns1="http://www.csapi.org/schema/parlayx/common/v2_1"><messageId>%s</messageId><text>%s</text><variables>%s</variables></ns1:ServiceException></detail></soapenv:Fault></soapenv:Body></soapenv:Envelope>`

var successResponse = `<?xml version="1.0" encoding="UTF-8"?><soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><soapenv:Body><ns1:sendSmsResponse xmlns:ns1="http://www.csapi.org/schema/parlayx/sms/send/v2_2/local"><ns1:result>%s</ns1:result></ns1:sendSmsResponse></soapenv:Body></soapenv:Envelope>`

// SMSEnvelope payload
type SMSEnvelope struct {
	SMSReqHeader struct {
		RequestHeader struct {
			SpID       string `xml:"spId"`
			SpPassword string `xml:"spPassword"`
			ServiceID  string `xml:"serviceId"`
			TimeStamp  string `xml:"timeStamp"`
			LinkID     string `xml:"linkid"`
			OA         string `xml:"OA"`
			FA         string `xml:"FA"`
		} `xml:"RequestSOAPHeader"`
	} `xml:"Header"`
	SMSReqBody struct {
		RequestBody struct {
			Number   string `xml:"addresses"`
			SenderID string `xml:"senderName"`
			Message  string `xml:"message"`
			RRequest struct {
				EndPoint   string `xml:"endpoint"`
				IntName    string `xml:"interfaceName"`
				Correlator string `xml:"correlator"`
			} `xml:"receiptRequest"`
		} `xml:"sendSms"`
	} `xml:"Body"`
}

// SafPage endpoint to rm request
func SafPage(w http.ResponseWriter, r *http.Request) {

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
		"SafRequest:: SenderID: %s, Phone: %s, Message: %s", reqBody.SenderID,
		reqBody.Number, reqBody.Message,
	))

	// senderIDs := []string{"FOCUSMOBILE", "Eutychus", "SMSLEOPARD", "601947"}

	// if !utils.InArray(reqBody.SenderID, senderIDs) {
	if !checkSID(reqBody.SenderID) {
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

func checkSID(sid string) bool {
	// var out = []bool{true, false}
	// return out[rand.Intn(len(out))]
	return true
}
