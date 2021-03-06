package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	sessionIdHeader = "X-Transmission-Session-Id"
	success         = "success"
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

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TransmissionClient struct {
	config     *TransmissionConfig
	httpClient HttpClient
	sessionId  string
}

func (client *TransmissionClient) GetTorrents() ([]string, error) {
	req := request{
		Method: "torrent-get",
		Arguments: map[string]interface{}{
			"fields": []string{"name"},
		},
	}
	res, _ := client.httpClient.Do(client.createHttpRequest(req))

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		panic(err)
	}
	result := responseBody.Result
	if result == success {
		torrents := responseBody.Arguments["torrents"].([]interface{})
		var list []string
		for _, torrent := range torrents {
			torrentName := (torrent.(map[string]interface{}))["name"].(string)
			list = append(list, torrentName)
		}
		return list, nil
	} else {
		return nil, errors.New(result)
	}
}

func (client *TransmissionClient) StopTorrents() error {
	req := request{
		Method: "torrent-stop",
	}
	res, _ := client.httpClient.Do(client.createHttpRequest(req))

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		panic(err)
	}
	result := responseBody.Result
	if result == success {
		return nil
	} else {
		return errors.New(result)
	}
}

func (client *TransmissionClient) StartTorrents() error {
	req := request{
		Method: "torrent-start",
	}
	rp, err := client.httpClient.Do(client.createHttpRequest(req))
	if err != nil {
		return err
	}
	defer rp.Body.Close()
	body, err := ioutil.ReadAll(rp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseBody response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		panic(err)
	}
	result := responseBody.Result
	if result == success {
		return nil
	} else {
		return errors.New(result)
	}
}

func (client *TransmissionClient) getSessionId() string {
	if client.sessionId == "" {
		log.Print("Session ID is empty")
		rq, _ := http.NewRequest(http.MethodPost, client.config.Url, nil)
		rq.SetBasicAuth(client.config.Username, client.config.Password)
		rp, err := client.httpClient.Do(rq)
		if err != nil {
			log.Fatalln(err)
		}

		defer rp.Body.Close()
		sessionId := rp.Header.Get(sessionIdHeader)
		client.sessionId = sessionId
		return sessionId
	} else {
		return client.sessionId
	}
}

func (client *TransmissionClient) createHttpRequest(req request) *http.Request {
	requestBody, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPost, client.config.Url, bytes.NewBuffer(requestBody))
	request.SetBasicAuth(client.config.Username, client.config.Password)
	request.Header.Set(sessionIdHeader, client.getSessionId())
	return request
}
