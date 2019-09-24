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
	Method    string `json:"method"`
	Arguments struct {
		Fields []string `json:"fields"`
	} `json:"arguments"`
}

type Response struct {
	Arguments struct {
		Torrents []struct {
			Name string `json:"name"`
		} `json:"torrents"`
	} `json:"arguments"`
	Result string `json:"result"`
}

func GetTorrents() []string {
	client := &http.Client{}
	req := Request{
		Method: "torrent-get",
		Arguments: struct {
			Fields []string `json:"fields"`
		}{
			Fields: []string{"name"},
		},
	}
	requestBody, err := json.Marshal(req)
	if err != nil {
		log.Fatalln(err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
	sessionID := response.Header.Get("X-Transmission-Session-Id")
	log.Println(sessionID)
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
	log.Println(string(body))
	var responseBody Response
	if err := json.Unmarshal([]byte(string(body)), &responseBody); err != nil {
		panic(err)
	}
	log.Println(responseBody)
	return []string{responseBody.Arguments.Torrents[0].Name}
}
