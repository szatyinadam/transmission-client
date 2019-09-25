package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type request struct {
	Method    string                 `json:"method" binding:"required"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Tag       int32                  `json:"tag,omitempty"`
}

type response struct {
	Result    string                 `json:"result" binding:"required"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Tag       int32                  `json:"tag,omitempty"`
}

func GetTorrents(transmission *Transmission) []string {
	client := &http.Client{}
	req := request{
		Method: "torrent-get",
		Arguments: map[string]interface{}{
			"fields": []string{"name"},
		},
	}
	requestBody, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", transmission.Url, bytes.NewBuffer(requestBody))
	res, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sessionID := res.Header.Get("X-Transmission-Session-Id")
	request.Header.Set("X-Transmission-Session-Id", sessionID)
	res, err = client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody response
	if err := json.Unmarshal([]byte(string(body)), &responseBody); err != nil {
		panic(err)
	}
	torrents := responseBody.Arguments["torrents"].([]interface{})
	var list []string
	for _, torrent := range torrents {
		torrentName := (torrent.(map[string]interface{}))["name"].(string)
		list = append(list, torrentName)
	}
	return list
}
