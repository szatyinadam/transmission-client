package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (httpClient *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if httpClient.DoFunc != nil {
		return httpClient.DoFunc(req)
	}
	return &http.Response{}, nil
}

func TestGetTorrents(t *testing.T) {
	json := `{
    	"arguments": {
        	"torrents": [
            	{
                	"name": "South Park S21 720p"
            	},
            	{
                	"name": "The.Curious.Case.of.Benjamin.Button.2008.iNT.720p.BluRay.DTS.x264.HuN-TRiNiTY"
            	}
        	]
    	},
    	"result": "success"
	}`
	client := getMockClient(json)

	torrents, _ := client.GetTorrents()

	assert.Equal(t, len(torrents), 2)
	assert.Equal(t, torrents[0], "South Park S21 720p")
	assert.Equal(t, torrents[1], "The.Curious.Case.of.Benjamin.Button.2008.iNT.720p.BluRay.DTS.x264.HuN-TRiNiTY")
}

func TestStopTorrents(t *testing.T) {
	json := `{"arguments": {}, "result": "success"}`
	client := getMockClient(json)

	assert.Nil(t, client.StopTorrents())
}

func TestStartTorrents(t *testing.T) {
	json := `{"arguments": {}, "result": "success"}`
	client := getMockClient(json)

	assert.Nil(t, client.StartTorrents())
}

func getMockClient(responseBody string) *TransmissionClient {
	reader := ioutil.NopCloser(bytes.NewReader([]byte(responseBody)))
	client := &TransmissionClient{
		config: &TransmissionConfig{
			Url:      "http://transmission.example.com:9091/transmission/rpc",
			Username: "transmission",
			Password: "password",
		},
		httpClient: &mockHttpClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       reader,
				}, nil
			},
		},
		sessionId: "Zpge0lNLRb8nx04peu4MgEuiWU6oY1ouFUI2RxiNmPCi38mm",
	}
	return client
}
