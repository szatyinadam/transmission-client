package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type TransmissionConfig struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SqsConfig struct {
	QueueUrl          string `yaml:"queue-url"`
	Region            string `yaml:"region"`
	CredentialPath    string `yaml:"credential-path"`
	CredentialProfile string `yaml:"credential-profile"`
	MaxMessages       int64  `yaml:"max-messages"`
}

type Config struct {
	Transmission TransmissionConfig `yaml:"transmission"`
	Sqs          SqsConfig          `yaml:"sqs"`
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
