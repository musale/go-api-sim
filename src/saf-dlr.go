package src

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var dlrResponse = `
<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
   <soapenv:Body>
      <ns1:getSmsDeliveryStatusResponse xmlns:ns1="http://www.csapi.org/schema/parlayx/sms/send/v2_2/local">
         <ns1:result>
            <address>tel:%s</address>
            <deliveryStatus>%s</deliveryStatus>
         </ns1:result>
      </ns1:getSmsDeliveryStatusResponse>
   </soapenv:Body>
</soapenv:Envelope>
`

type DLREnvelope struct {
	Header struct {
		SOAPHeader struct {
			SPID       string `xml:"spId"`
			SPPAssword string `xml:"spPassword"`
			ServiceID  string `xml:"serviceId"`
			TimeStamp  string `xml:"timeStamp"`
			OA         string `xml:"OA"`
			FA         string `xml:"FA"`
		} `xml:"RequestSOAPHeader"`
	} `xml:"Header"`
	Body struct {
		DLRStatus struct {
			RequestID string `xml:"requestIdentifier"`
		} `xml:"getSmsDeliveryStatus"`
	} `xml:"Body"`
}

// SafDLRPage endpoint to dlr status query
func SafDLRPage(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ioutil readall: ", err)
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.Write([]byte("read all error"))
		return
	}

	var req DLREnvelope
	if err := xml.Unmarshal(body, &req); err != nil {
		log.Println("Xml unmarshal: ", err)
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.Write([]byte("xml unmarshal err"))
		return
	}

	reqID := req.Body.DLRStatus.RequestID
	number := req.Header.SOAPHeader.OA
	if number[0:4] == "tel:" {
		number = number[4:]
	}
	log.Println(fmt.Sprintf(
		"Get status for ID: %s, Phone: %s", reqID, number))

	status := "DeliveredToTerminal"
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
	w.Write([]byte(fmt.Sprintf(dlrResponse, number, status)))
	return
}
