package src

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/etowett/go-api-sim/utils"
)

var senderIDResponse = `
<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
   <soapenv:Body>
      <soapenv:Fault>
         <faultcode>%s</faultcode>
         <faultstring>%s</faultstring>
         <detail>
            <ns1:ServiceException xmlns:ns1="http://www.csapi.org/schema/parlayx/common/v2_1">
               <messageId>%s</messageId>
               <text>SenderName or senderAddress is unknown!</text>
               <variables>FOCUSMOBILE</variables>
            </ns1:ServiceException>
         </detail>
      </soapenv:Fault>
   </soapenv:Body>
</soapenv:Envelope>
`

var successResponse = `
<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
   <soapenv:Body>
      <ns1:sendSmsResponse xmlns:ns1="http://www.csapi.org/schema/parlayx/sms/send/v2_2/local">
         <ns1:result>%s</ns1:result>
      </ns1:sendSmsResponse>
   </soapenv:Body>
</soapenv:Envelope>
`

// SMSEnvelope payload
type SMSEnvelope struct {
	XMLName xml.Name
	Header  SMSHeader `xml:"Header"`
	ReqBody SMSBody   `xml:"Body"`
}

// SMSHeader Payload header
type SMSHeader struct {
	XMLName    xml.Name `xml:"RequestSOAPHeader"`
	spID       string   `xml:"spId"`
	spPassword string   `xml:"spPassword"`
	serviceID  string   `xml:"serviceId"`
	timeStamp  string   `xml:"timeStamp"`
	linkID     string   `xml:"linkid"`
	OA         string   `xml:"OA"`
	FA         string   `xml:"FA"`
}

// SMSBody payload body
type SMSBody struct {
	XMLName  xml.Name `xml:"sendSms"`
	Number   string   `xml:"addresses"`
	SenderID string   `xml:"senderName"`
	Message  string   `xml:"message"`
	RRequest RRequest `xml:"receiptRequest"`
}

// RRequest recipient request
type RRequest struct {
	EndPoint   string `xml:"endpoint"`
	IntName    string `xml:"interfaceName"`
	Correlator string `xml:"correlator"`
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
	log.Println("received: ", body)

	var req SMSEnvelope
	if err := xml.Unmarshal(body, &req); err != nil {
		log.Println("Xml unmarshal: ", err)
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.Write([]byte("xml unmarshal err"))
		return
	}

	reqBody := req.ReqBody
	load := fmt.Sprintf(
		"SenderID: %s, Phone: %s, Message: %s", reqBody.SenderID,
		reqBody.Number, reqBody.Message,
	)
	log.Println(load)

	senderIDs := []string{"FOCUSMOBILE", "Eutychus", "SMSLEOPARD"}

	if utils.InArray(reqBody.SenderID, senderIDs) {
		faultCode := "SVC0002"
		faultString := "SenderName or senderAddress is unknown!"
		retResponse := fmt.Sprintf(
			senderIDResponse, faultCode, faultString, faultCode, faultString,
			req.ReqBody.SenderID,
		)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.Write([]byte(retResponse))
		return
	}
	msgID := "100001200501170419072620015931"
	retResponse := fmt.Sprintf(successResponse, msgID)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
	w.Write([]byte(retResponse))
	return
}
