package main

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	test := struct {
		name string
		want *Config
	}{
		name: "Load config",
		want: &Config{
			Transmission: Transmission{
				Url: "http://192.168.0.238:9091/transmission/rpc",
			},
			Sqs: Sqs{
				QueueUrl:          "https://sqs.us-east-1.amazonaws.com/356329984695/transmission.fifo",
				Region:            "us-east-1",
				CredentialPath:    "/Users/adam/.aws/credentials",
				CredentialProfile: "default",
			},
		},
	}
	t.Run(test.name, func(t *testing.T) {
		if got := ReadConfig(); !reflect.DeepEqual(got, test.want) {
			t.Errorf("ReadConfig() = %v, want %v", got, test.want)
		}
	})
}
