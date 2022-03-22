package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	type Body struct {
		Id    int
		Phone int
		Text  string
	}
	debugSpam := &Body{
		Id:    11,
		Phone: 72222222222,
		Text:  "ASDASDASDASDasdsadsadssssadssssssasdasdsstsasdsstsssddsdASDASDasdasdsad",
	}
	messageIDdebug := 11
	url := "https://probe.fbrq.cloud/v1/send/"
	var bearer = "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzkwNjMwNjgsImlzcyI6ImZhYnJpcXVlIiwibmFtZSI6IkVnb3IifQ.Ns0At0MuQyUqNpeZk5qbxdhjzKC5QQ_NVLnhoqvQxc0"

	jsonBody, err := json.Marshal(debugSpam)
	if err != nil {
		return
	}

	fmt.Println("json Body", bytes.NewBuffer(jsonBody))
	fmt.Println(url + strconv.Itoa(messageIDdebug))
	req, err := http.NewRequest("POST", url+strconv.Itoa(messageIDdebug), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println(req.Header)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string([]byte(body)))
}
