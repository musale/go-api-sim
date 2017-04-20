package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/etowett/garbanzo/utils"
)

const (
	SMS_URL   = "http://goapi.local/saf"
	SMS_QUERY = `
<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:loc="http://www.csapi.org/schema/parlayx/sms/send/v2_2/local" xmlns:v2="http://www.huawei.com.cn/schema/common/v2_1">
    <soapenv:Header>
        <v2:RequestSOAPHeader>
            <v2:spId>%s</v2:spId>
            <v2:spPassword>%s</v2:spPassword>
            <v2:serviceId>%s</v2:serviceId>
            <v2:timeStamp>%s</v2:timeStamp>
            <v2:linkid>%s</v2:linkid>
            <v2:OA>tel:%s</v2:OA>
            <v2:FA>tel:%s</v2:FA>
        </v2:RequestSOAPHeader>
    </soapenv:Header>
    <soapenv:Body>
        <loc:sendSms>
            <loc:addresses>tel:%s</loc:addresses>
            <loc:senderName>%s</loc:senderName>
            <loc:message>%s</loc:message>
            <loc:receiptRequest>
                <endpoint>http://10.138.30.123:9080/notify</endpoint>
                <interfaceName>SmsNotification</interfaceName>
                <correlator>123</correlator>
            </loc:receiptRequest>
        </loc:sendSms>
    </soapenv:Body>
</soapenv:Envelope>
`
)

// Response return
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// SuccessEnvelope payload body
type SuccessEnvelope struct {
	SuccessBody struct {
		SuccessResponse struct {
			Result string `xml:"result"`
		} `xml:"sendSmsResponse"`
	} `xml:"Body"`
}

// SIDEnvelope payload body
type SIDEnvelope struct {
	SIDBody struct {
		SIDFault struct {
			FaultCode   string `xml:"faultcode"`
			FaultString string `xml:"faultstring"`
			Detail      struct {
				ServiceException struct {
					MessageID string `xml:"messageId"`
					Text      string `xml:"text"`
					Variables string `xml:"variables"`
				} `xml:"ServiceException"`
			} `xml:"detail"`
		} `xml:"Fault"`
	} `xml:"Body"`
}

func main() {

	spID := "123123"
	spPassword := "Scre!sdfsg"
	timeStamp := time.Now().Format("20060102150405")
	passString := spID + spPassword + timeStamp
	realPass := utils.GetMD5Hash(passString)
	serviceID := "122122122122122"
	linkID := "123456"

	senderName := "HolyHigh"
	// senderName := "FOCUSMOBILE"
	destNum := "254715458745"
	message := "Hello world!"

	requestContent := fmt.Sprintf(
		SMS_QUERY, spID, realPass, serviceID, timeStamp, linkID, destNum,
		destNum, destNum, senderName, message,
	)

	httpClient := new(http.Client)
	resp, err := httpClient.Post(
		SMS_URL, "text/xml; charset=utf-8",
		bytes.NewBufferString(requestContent),
	)
	if err != nil {
		log.Println("Post error: ", err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("read all error: ", err)
	}

	var env SuccessEnvelope
	if err := xml.Unmarshal(b, &env); err != nil {
		log.Println("success unmarshal error: ", err)
		return
	}
	result := env.SuccessBody.SuccessResponse.Result
	if len(result) < 1 {
		var sid SIDEnvelope
		if err := xml.Unmarshal(b, &sid); err != nil {
			log.Println("sid unmarshal error: ", err)
			return
		}
		log.Println("Response: ", sid)
	} else {
		log.Println("Response: ", env)
	}
	return
}
