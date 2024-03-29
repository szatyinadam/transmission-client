package main

import (
	"github.com/stretchr/testify/assert"
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
			Transmission: TransmissionConfig{
				Url:      "http://192.168.0.238:9091/transmission/rpc",
				Username: "transmission",
				Password: "password",
			},
			Sqs: SqsConfig{
				QueueUrl:          "https://sqs.us-east-1.amazonaws.com/356329984695/transmission.fifo",
				Region:            "us-east-1",
				CredentialPath:    "/Users/adam/.aws/credentials",
				CredentialProfile: "default",
				MaxMessages:       5,
			},
		},
	}
	t.Run(test.name, func(t *testing.T) {
		config := ReadConfig()
		assert.True(t, reflect.DeepEqual(config, test.want))
	})
}
