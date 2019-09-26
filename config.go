package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Transmission struct {
	Url string `yaml:"url"`
}

type Sqs struct {
	QueueUrl          string `yaml:"queue-url"`
	Region            string `yaml:"region"`
	CredentialPath    string `yaml:"credential-path"`
	CredentialProfile string `yaml:"credential-profile"`
}

type Config struct {
	Transmission Transmission `yaml:"transmission"`
	Sqs          Sqs          `yaml:"sqs"`
}

func ReadConfig() *Config {
	f, err := os.Open("config/config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	decoder := yaml.NewDecoder(f)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}
