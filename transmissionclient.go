package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "http://192.168.0.238:9091/transmission/rpc"

type Request struct {
	Method    string                 `json:"method" binding:"required"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Tag       int32                  `json:"tag,omitempty"`
}

type Response struct {
	Result    string                 `json:"result" binding:"required"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Tag       int32                  `json:"tag,omitempty"`
}

func GetTorrents() []string {
	client := &http.Client{}
	req := Request{
		Method: "torrent-get",
		Arguments: map[string]interface{}{
			"fields": []string{"name"},
		},
	}
	requestBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sessionID := response.Header.Get("X-Transmission-Session-Id")
	request.Header.Set("X-Transmission-Session-Id", sessionID)
	response, err = client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody Response
	if err := json.Unmarshal([]byte(string(body)), &responseBody); err != nil {
		panic(err)
	}
	log.Println(responseBody)
	return []string{}
}
