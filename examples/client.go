package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	aftURL = "http://127.0.0.1:4027/aft"
	rmsURL = "http://127.0.0.1:4027/routesms"
)

func main() {

	client := http.Client{}
	form := url.Values{}
	form.Add("username", "etowett")
	form.Add("message", "Hello world")
	form.Add("to", "727328542")
	form.Add("from", "Eutychus")
	req, err := http.NewRequest(
		"POST", aftURL, strings.NewReader(form.Encode()))

	if err != nil {
		log.Fatal("Request Error: ", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Accept", "Application/json")
	req.Header.Add("apikey", "1234#")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Do error: ", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("ReadAll error: ", err)
	}

	log.Println("AFT: ", string(body))
	return
}

func getAccount() string {
	return "admin"
}

func getAmount() string {
	return "123.45"
}

func getPhone() string {
	return "254727372285"
}

func getTransTime() string {
	return time.Now().String()
}
