package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Transmission struct {
	Url string `yaml:"url"`
}

type Config struct {
	Transmission Transmission `yaml:"transmission"`
}

func ReadConfig(config *Config) {
	f, err := os.Open("config/config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		log.Fatalln(err)
	}
}
